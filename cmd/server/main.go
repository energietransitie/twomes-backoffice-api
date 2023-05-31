package main

import (
	"context"
	"io/fs"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/handlers"
	"github.com/energietransitie/twomes-backoffice-api/repositories"
	"github.com/energietransitie/twomes-backoffice-api/services"
	"github.com/energietransitie/twomes-backoffice-api/swaggerdocs"
	"github.com/energietransitie/twomes-backoffice-api/twomes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

// Configuration holds all the configuration for the server.
type Configuration struct {
	DatabaseDSN string
	BaseURL     string
}

func getConfiguration() Configuration {
	dsn, ok := os.LookupEnv("TWOMES_DSN")
	if !ok {
		logrus.Fatal("TWOMES_DSN was not set")
	}

	baseURL, ok := os.LookupEnv("TWOMES_BASE_URL")
	if !ok {
		logrus.Fatal("TWOMES_BASE_URL was not set")
	}

	return Configuration{
		DatabaseDSN: dsn,
		BaseURL:     baseURL,
	}
}

func main() {
	config := getConfiguration()

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := repositories.NewDatabaseConnectionAndMigrate(ctx, config.DatabaseDSN)
	if err != nil {
		logrus.Fatal(err)
	}

	authService, err := services.NewAuthorizationServiceFromFile("./data/key.pem")
	if err != nil {
		logrus.Fatal(err)
	}
	authHandler := handlers.NewAuthorizationHandler(authService)

	adminRepository, err := repositories.NewAdminRepository("./data/admins.db")
	if err != nil {
		logrus.Fatal(err)
	}
	adminService := services.NewAdminService(adminRepository, authService)
	adminHandler := handlers.NewAdminHandler(adminService)

	adminAuth := authHandler.Middleware(twomes.AdminToken)
	accountActivationAuth := authHandler.Middleware(twomes.AccountActivationToken)
	accountAuth := authHandler.Middleware(twomes.AccountToken)
	deviceAuth := authHandler.Middleware(twomes.DeviceToken)

	appRepository := repositories.NewAppRepository(db)
	cloudFeedRepository := repositories.NewCloudFeedRepository(db)
	cloudFeedAuthRepository := repositories.NewCloudFeedAuthRepository(db)
	campaignRepository := repositories.NewCampaignRepository(db)
	buildingRepository := repositories.NewBuildingRepository(db)
	accountRepository := repositories.NewAccountRepository(db)
	propertyRepository := repositories.NewPropertyRepository(db)
	deviceTypeRepository := repositories.NewDeviceTypeRepository(db)
	deviceRepository := repositories.NewDeviceRepository(db)
	uploadRepository := repositories.NewUploadRepository(db)

	appService := services.NewAppService(appRepository)
	cloudFeedService := services.NewCloudFeedService(cloudFeedRepository)
	cloudFeedAuthService := services.NewCloudFeedAuthService(cloudFeedAuthRepository)
	campaignService := services.NewCampaignService(campaignRepository, appService, cloudFeedService)
	buildingService := services.NewBuildingService(buildingRepository)
	accountService := services.NewAccountService(accountRepository, authService, appService, campaignService, buildingService)
	propertyService := services.NewPropertyService(propertyRepository)
	deviceTypeService := services.NewDeviceTypeService(deviceTypeRepository, propertyService)
	deviceService := services.NewDeviceService(deviceRepository, authService, deviceTypeService, buildingService)
	uploadService := services.NewUploadService(uploadRepository, propertyService)

	appHandler := handlers.NewAppHandler(appService)
	cloudFeedHandler := handlers.NewCloudFeedHandler(cloudFeedService)
	cloudFeedAuthHandler := handlers.NewCloudFeedAuthHandler(cloudFeedAuthService)
	campaignHandler := handlers.NewCampaignHandler(campaignService)
	buildingHandler := handlers.NewBuildingHandler(buildingService)
	accountHandler := handlers.NewAccountHandler(accountService)
	deviceTypeHandler := handlers.NewDeviceTypeHandler(deviceTypeService)
	deviceHandler := handlers.NewDeviceHandler(deviceService)
	uploadHandler := handlers.NewUploadHandler(uploadService)

	r := chi.NewRouter()
	r.Use(middleware.Timeout(time.Second * 30))
	r.Use(middleware.Heartbeat("/healthcheck")) // Endpoint for health check.
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logrus.StandardLogger()}))

	r.Method("POST", "/app", adminAuth(adminHandler.Middleware(appHandler.Create))) // POST on /app.

	r.Method("POST", "/cloud_feed", adminAuth(adminHandler.Middleware(cloudFeedHandler.Create))) // POST on /cloud_feed.

	r.Method("POST", "/cloud_feed_auth", accountAuth(cloudFeedAuthHandler.Create)) // POST on /cloud_feed_auth.

	r.Method("POST", "/campaign", adminAuth(adminHandler.Middleware(campaignHandler.Create))) // POST on /campaign.

	r.Route("/account", func(r chi.Router) {
		r.Method("POST", "/", adminAuth(adminHandler.Middleware(accountHandler.Create))) // POST on /account.
		r.Method("POST", "/activate", accountActivationAuth(accountHandler.Activate))    // POST on /account/activate.
		r.Method("GET", "/{account_id}", accountAuth(accountHandler.GetAccountByID))     // GET on /account/{account_id}.
	})

	r.Method("GET", "/building/{building_id}", accountAuth(buildingHandler.GetBuildingByID)) // GET on /building/{building_id}.

	r.Method("POST", "/device_type", adminAuth(adminHandler.Middleware(deviceTypeHandler.Create))) // POST on /device_type.

	r.Route("/device", func(r chi.Router) {
		r.Method("POST", "/", accountAuth(deviceHandler.Create))                      // POST on /device.
		r.Method("POST", "/activate", handlers.Handler(deviceHandler.Activate))       // POST on /device/activate.
		r.Method("GET", "/{device_name}", accountAuth(deviceHandler.GetDeviceByName)) // GET on /device/{device_name}.
	})

	r.Method("POST", "/upload", deviceAuth(uploadHandler.Create)) // POST on /upload.

	setupSwaggerDocs(r, config.BaseURL)

	go setupAdminRPCHandler(adminHandler)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		logrus.Fatal(err)
	}
}

func setupSwaggerDocs(r *chi.Mux, baseURL string) {
	swaggerUI, err := fs.Sub(swaggerdocs.StaticFiles, "swagger-ui")
	if err != nil {
		logrus.Fatal(err)
	}

	docsHandler, err := handlers.NewDocsHandler(swaggerdocs.StaticFiles, baseURL)
	if err != nil {
		logrus.Fatal(err)
	}

	r.Method("GET", "/openapi.yml", handlers.Handler(docsHandler.OpenAPISpec))                        // Serve openapi.yml
	r.Method("GET", "/docs/*", http.StripPrefix("/docs/", http.FileServer(http.FS(swaggerUI))))       // Serve static files.
	r.Method("GET", "/docs", handlers.Handler(docsHandler.RedirectDocs(http.StatusMovedPermanently))) // Redirect /docs to /docs/
	r.Method("GET", "/", handlers.Handler(docsHandler.RedirectDocs(http.StatusSeeOther)))             // Redirect / to /docs/
}

func setupAdminRPCHandler(adminHandler *handlers.AdminHandler) {

	rpc.Register(adminHandler)
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp4", "127.0.0.1:8081")
	if err != nil {
		logrus.Fatal(err)
	}

	err = http.Serve(listener, nil)
	if err != nil {
		logrus.Fatal(err)
	}
}
