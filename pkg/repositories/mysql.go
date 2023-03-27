package repositories

import (
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
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
func NewDatabaseConnectionAndMigrate(dsn string) (*gorm.DB, error) {
	db, err := NewDatabaseConnection(dsn)
	if err != nil {
		return nil, err
	}

	return db, db.AutoMigrate(&AppModel{}, &CampaignModel{}, &AccountModel{}, &BuildingModel{}, &PropertyModel{}, &UploadModel{}, &DeviceTypeModel{}, &DeviceModel{}, &MeasurementModel{})
}
