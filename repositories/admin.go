package repositories

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AdminRepository struct {
	db *gorm.DB
}

// Create a new AdminRepository from a badger DB at fileName.
func NewAdminRepository(fileName string) (*AdminRepository, error) {
	db, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return &AdminRepository{}, err
	}

	err = db.AutoMigrate(&AdminModel{})
	if err != nil {
		return &AdminRepository{}, err
	}

	return &AdminRepository{
		db: db,
	}, nil
}

// Database representation of a [twomes.Admin].
type AdminModel struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	ActivatedAt time.Time
	Expiry      time.Time
}

// Set the name of the table in the database.
func (AdminModel) TableName() string {
	return "admin"
}

// Create a new AdminModel from a [twomes.Admin].
func MakeAdminModel(admin twomes.Admin) AdminModel {
	return AdminModel{
		Model:       gorm.Model{ID: admin.ID},
		Name:        admin.Name,
		ActivatedAt: admin.ActivatedAt,
		Expiry:      admin.Expiry,
	}
}

// Create a [twomes.Admin] from an AdminModel.
func (m *AdminModel) fromModel() twomes.Admin {
	return twomes.Admin{
		ID:          m.Model.ID,
		Name:        m.Name,
		ActivatedAt: m.ActivatedAt,
		Expiry:      m.Expiry,
	}
}

func (r *AdminRepository) Find(admin twomes.Admin) (twomes.Admin, error) {
	adminModel := MakeAdminModel(admin)
	err := r.db.Where(&adminModel).First(&adminModel).Error
	return adminModel.fromModel(), err
}

func (r *AdminRepository) GetAll() ([]twomes.Admin, error) {
	admins := make([]twomes.Admin, 0)

	var adminModels []AdminModel
	err := r.db.Find(&adminModels).Error
	if err != nil {
		return nil, err
	}

	for _, adminModel := range adminModels {
		admins = append(admins, adminModel.fromModel())
	}

	return admins, nil
}

func (r *AdminRepository) Create(admin twomes.Admin) (twomes.Admin, error) {
	adminModel := MakeAdminModel(admin)
	err := r.db.Create(&adminModel).Error
	return adminModel.fromModel(), err
}

func (r *AdminRepository) Update(admin twomes.Admin) (twomes.Admin, error) {
	adminModel := MakeAdminModel(admin)
	err := r.db.Model(&adminModel).Updates(adminModel).Error
	return adminModel.fromModel(), err
}

func (r *AdminRepository) Delete(admin twomes.Admin) error {
	adminModel := MakeAdminModel(admin)
	return r.db.Where(&adminModel).Delete(&adminModel).Error
}
