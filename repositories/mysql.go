package repositories

import (
	"context"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Create a new database connection that can be used by repositories.
func NewDatabaseConnection(dsn string) (*gorm.DB, error) {
	cfg, err := mysqldriver.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	cfg.ParseTime = true
	cfg.Loc = time.UTC
	cfg.Params = map[string]string{"charset": "utf8mb4"}

	dsn = cfg.FormatDSN()

	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

// Create a new database connection and perform a migration.
// Try until connection succeeds or context is done.
func NewDatabaseConnectionAndMigrate(ctx context.Context, dsn string) (db *gorm.DB, err error) {
	_, ok := ctx.Deadline()
	if !ok {
		logrus.Warn("no deadline was set for making database connection. we will try indefinately")
	}

	logrus.Info("connecting to database")

	for {
		db, err = NewDatabaseConnection(dsn)
		if err == nil {
			return db, db.AutoMigrate(&AppModel{}, &CloudFeedModel{}, &CampaignModel{}, &AccountModel{}, &CloudFeedAuthModel{}, &BuildingModel{}, &PropertyModel{}, &UploadModel{}, &DeviceTypeModel{}, &DeviceModel{}, &MeasurementModel{})
		}

		select {
		case <-time.After(time.Second): // Wait for 1 second before we loop again.
		case <-ctx.Done():
			return // Return with the database and error.
		}
	}
}
