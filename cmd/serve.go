package cmd

import (
	"context"
	"io/fs"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/energietransitie/needforheat-server-api/handlers"
	"github.com/energietransitie/needforheat-server-api/needforheat/authorization"
	"github.com/energietransitie/needforheat-server-api/repositories"
	"github.com/energietransitie/needforheat-server-api/services"
	"github.com/energietransitie/needforheat-server-api/swaggerdocs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the needforheat API server",
		RunE:  handleServe,
	}

	rootCmd.AddCommand(serveCmd)
}

const (
	shutdownTimeout    = 30 * time.Second
	preRenewalDuration = 12 * time.Hour
)

func handleServe(cmd *cobra.Command, args []string) error {
	config := getConfiguration()

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	dbCtx, dbCancel := context.WithTimeout(ctx, 10*time.Second)
	defer dbCancel()

	db, err := repositories.NewDatabaseConnectionAndMigrate(dbCtx, config.DatabaseDSN)

	if err != nil {
		logrus.Fatal(err)
	}

	//Important services for admin and auth
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

	adminAuth := authHandler.Middleware(authorization.AdminToken)
	accountActivationAuth := authHandler.Middleware(authorization.AccountActivationToken)
	accountAuth := authHandler.Middleware(authorization.AccountToken)
	deviceAuth := authHandler.Middleware(authorization.DeviceToken)

	//Repositories
	appRepository := repositories.NewAppRepository(db)
	cloudFeedTypeRepository := repositories.NewCloudFeedTypeRepository(db)
	cloudFeedRepository := repositories.NewCloudFeedRepository(db)
	campaignRepository := repositories.NewCampaignRepository(db)
	propertyRepository := repositories.NewPropertyRepository(db)
	uploadRepository := repositories.NewUploadRepository(db)
	buildingRepository := repositories.NewBuildingRepository(db)
	accountRepository := repositories.NewAccountRepository(db)
	deviceTypeRepository := repositories.NewDeviceTypeRepository(db)
	deviceRepository := repositories.NewDeviceRepository(db)
	dataSourceListRepository := repositories.NewDataSourceListRepository(db)
	dataSourceTypeRepository := repositories.NewDataSourceTypeRepository(db)

	//Services
	appService := services.NewAppService(appRepository)
	cloudFeedTypeService := services.NewCloudFeedTypeService(cloudFeedTypeRepository)
	propertyService := services.NewPropertyService(propertyRepository)
	deviceTypeService := services.NewDeviceTypeService(deviceTypeRepository, propertyService)
	dataSourceTypeService := services.NewDataSourceTypeService(
		dataSourceTypeRepository,
		deviceTypeService,
		cloudFeedTypeService,
	)
	dataSourceListService := services.NewDataSourceListService(dataSourceListRepository, dataSourceTypeService)
	campaignService := services.NewCampaignService(campaignRepository, appService, dataSourceListService)
	uploadService := services.NewUploadService(uploadRepository, deviceRepository, propertyService)
	cloudFeedService := services.NewCloudFeedService(cloudFeedRepository, cloudFeedTypeRepository, uploadService)
	buildingService := services.NewBuildingService(buildingRepository, uploadService)
	accountService := services.NewAccountService(accountRepository, authService, appService, campaignService, buildingService, cloudFeedService, dataSourceTypeService)
	deviceService := services.NewDeviceService(deviceRepository, authService, deviceTypeService, buildingService, uploadService)

	//Handlers
	appHandler := handlers.NewAppHandler(appService)
	cloudFeedTypeHandler := handlers.NewCloudFeedTypeHandler(cloudFeedTypeService)
	cloudFeedHandler := handlers.NewCloudFeedHandler(cloudFeedService)
	campaignHandler := handlers.NewCampaignHandler(campaignService)
	uploadHandler := handlers.NewUploadHandler(uploadService)
	buildingHandler := handlers.NewBuildingHandler(buildingService)
	accountHandler := handlers.NewAccountHandler(accountService)
	deviceTypeHandler := handlers.NewDeviceTypeHandler(deviceTypeService)
	deviceHandler := handlers.NewDeviceHandler(deviceService)
	dataSourceListHandler := handlers.NewDataSourceListHandler(dataSourceListService)
	dataSourceTypeHandler := handlers.NewDataSourceTypeHandler(dataSourceTypeService)

	go cloudFeedService.RefreshTokensInBackground(ctx, preRenewalDuration)
	go cloudFeedService.DownloadInBackground(ctx, config.downloadStartTime)

	//Router
	r := chi.NewRouter()

	r.Use(middleware.Timeout(time.Second * 30))
	r.Use(middleware.Heartbeat("/healthcheck")) // Endpoint for health check.
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logrus.StandardLogger()}))

	r.Method("POST", "/app", adminAuth(adminHandler.Middleware(appHandler.Create))) // POST on /app.

	r.Method("POST", "/cloud_feed_type", adminAuth(adminHandler.Middleware(cloudFeedTypeHandler.Create))) // POST on /cloud_feed.

	r.Method("POST", "/campaign", adminAuth(adminHandler.Middleware(campaignHandler.Create))) // POST on /campaign.

	r.Route("/account", func(r chi.Router) {
		r.Method("POST", "/", adminAuth(adminHandler.Middleware(accountHandler.Create))) // POST on /account.
		r.Method("POST", "/activate", accountActivationAuth(accountHandler.Activate))    // POST on /account/activate.

		r.Route("/{account_id}", func(r chi.Router) {
			r.Method("GET", "/", accountAuth(accountHandler.GetAccountByID))                     // GET on /account/{account_id}.
			r.Method("POST", "/cloud_feed", accountAuth(cloudFeedHandler.Create))                // POST on /account/{account_id}/cloud_feed_auth.
			r.Method("GET", "/cloud_feed", accountAuth(accountHandler.GetCloudFeedAuthStatuses)) // GET on /account/{account_id}/cloud_feed_auth.
		})
	})

	r.Method("GET", "/building/{building_id}", accountAuth(buildingHandler.GetBuildingByID)) // GET on /building/{building_id}.

	r.Method("POST", "/device_type", adminAuth(adminHandler.Middleware(deviceTypeHandler.Create))) // POST on /device_type.

	r.Route("/device", func(r chi.Router) {
		r.Method("POST", "/", accountAuth(deviceHandler.Create))                                         // POST on /device.
		r.Method("POST", "/activate", handlers.Handler(deviceHandler.Activate))                          // POST on /device/activate.
		r.Method("GET", "/{device_name}", accountAuth(deviceHandler.GetDeviceByName))                    // GET on /device/{device_name}.
		r.Method("GET", "/{device_name}/measurements", accountAuth(deviceHandler.GetDeviceMeasurements)) // GET on /device/{device_name}/measurements.
		r.Method("GET", "/{device_name}/properties", accountAuth(deviceHandler.GetDeviceProperties))     // GET on /device/{device_name}/properties.
	})

	r.Method("POST", "/upload", deviceAuth(uploadHandler.Create)) // POST on /upload.

	r.Method("POST", "/data_source_list", adminAuth(dataSourceListHandler.Create)) // POST on /data_source_list
	r.Method("POST", "/data_source_type", adminAuth(dataSourceTypeHandler.Create)) // POST on /data_source_type

	setupSwaggerDocs(r, config.BaseURL)

	go setupRPCHandler(adminHandler, cloudFeedHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	err = listenAndServe(ctx, server)
	if err != nil {
		return err
	}
	return nil
}

const (
	day                 = time.Hour * 24
	defaultDownloadTime = "04h00s"
)

type Configuration struct {
	DatabaseDSN       string
	BaseURL           string
	downloadStartTime time.Time
}

func getConfiguration() Configuration {
	dsn, ok := os.LookupEnv("NFH_DSN")
	if !ok {
		logrus.Fatal("NFH_DSN was not set")
	}

	baseURL, ok := os.LookupEnv("NFH_BASE_URL")
	if !ok {
		logrus.Fatal("NFH_BASE_URL was not set")
	}

	downloadTime, ok := os.LookupEnv("NFH_DOWNLOAD_TIME")
	if !ok {
		logrus.Warning("NFH_DOWNLOAD_TIME was not set. defaulting to", defaultDownloadTime)
		downloadTime = defaultDownloadTime
	}

	duration, err := time.ParseDuration(downloadTime)
	if err != nil {
		logrus.Fatal(err)
	}

	downloadStartTime := time.Now().Truncate(day)
	downloadStartTime = downloadStartTime.Add(duration)
	// If time is in the past, add 1 day.
	if downloadStartTime.Before(time.Now()) {
		downloadStartTime = downloadStartTime.Add(day)
	}

	return Configuration{
		DatabaseDSN:       dsn,
		BaseURL:           baseURL,
		downloadStartTime: downloadStartTime,
	}
}

func listenAndServe(ctx context.Context, server *http.Server) error {
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		err := server.ListenAndServe()
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	})
	logrus.Infoln("listening on", server.Addr)

	g.Go(func() error {
		<-gCtx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return server.Shutdown(shutdownCtx)
	})

	err := g.Wait()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
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

func setupRPCHandler(handlers ...any) {

	for _, handler := range handlers {
		rpc.Register(handler)
	}

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
