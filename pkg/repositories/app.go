package repositories

import (
	"github.com/energietransitie/twomes-api/pkg/twomes"
	"gorm.io/gorm"
)

type AppRepository struct {
	db *gorm.DB
}

func NewAppRepository(db *gorm.DB) *AppRepository {
	return &AppRepository{
		db: db,
	}
}

// Database representation of a [twomes.App].
type AppModel struct {
	gorm.Model
	Name                    string `gorm:"unique;not null"`
	ProvisioningURLTemplate string
}

// Set the name of the table in the database.
func (AppModel) TableName() string {
	return "apps"
}

// Create a new AppModel from a [twomes.App]
func MakeAppModel(app twomes.App) AppModel {
	return AppModel{
		Model:                   gorm.Model{ID: app.ID},
		Name:                    app.Name,
		ProvisioningURLTemplate: app.ProvisioningURLTemplate,
	}
}

// Create a [twomes.App] from an AppModel.
func (m *AppModel) fromModel() twomes.App {
	return twomes.App{
		ID:                      m.Model.ID,
		Name:                    m.Name,
		ProvisioningURLTemplate: m.ProvisioningURLTemplate,
	}
}

func (r *AppRepository) Find(app twomes.App) (twomes.App, error) {
	appModel := MakeAppModel(app)
	err := r.db.Where(&appModel).First(&appModel).Error
	return appModel.fromModel(), err
}

func (r *AppRepository) GetAll() ([]twomes.App, error) {
	apps := make([]twomes.App, 0)

	var appModels []AppModel
	err := r.db.Find(&appModels).Error
	if err != nil {
		return nil, err
	}

	for _, appModel := range appModels {
		apps = append(apps, appModel.fromModel())
	}

	return apps, nil
}

func (r *AppRepository) Create(app twomes.App) (twomes.App, error) {
	appModel := MakeAppModel(app)
	err := r.db.Create(&appModel).Error
	return appModel.fromModel(), err
}

func (r *AppRepository) Delete(app twomes.App) error {
	appModel := MakeAppModel(app)
	return r.db.Delete(&appModel).Error
}
