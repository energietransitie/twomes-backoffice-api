package repositories

import (
	"github.com/energietransitie/twomes-backoffice-api/pkg/twomes"
	"gorm.io/gorm"
)

type PropertyRepository struct {
	db *gorm.DB
}

// Create a new PropertyRepository.
func NewPropertyRepository(db *gorm.DB) *PropertyRepository {
	return &PropertyRepository{
		db: db,
	}
}

// Database representation of a [twomes.Property]
type PropertyModel struct {
	gorm.Model
	Name string `gorm:"unique;non null"`
}

// Set the name of the table in the database.
func (PropertyModel) TableName() string {
	return "property"
}

// Create a PropertyModel from a [twomes.Property].
func MakePropertyModel(property twomes.Property) PropertyModel {
	return PropertyModel{
		Model: gorm.Model{ID: property.ID},
		Name:  property.Name,
	}
}

// Create a [twomes.Property] from a PropertyModel.
func (m *PropertyModel) fromModel() twomes.Property {
	return twomes.Property{
		ID:   m.Model.ID,
		Name: m.Name,
	}
}

func (r *PropertyRepository) Find(property twomes.Property) (twomes.Property, error) {
	propertyModel := MakePropertyModel(property)
	err := r.db.Where(&propertyModel).First(&propertyModel).Error
	return propertyModel.fromModel(), err
}

func (r *PropertyRepository) GetAll() ([]twomes.Property, error) {
	var properties []twomes.Property

	var propertyModels []PropertyModel
	err := r.db.Find(&propertyModels).Error
	if err != nil {
		return nil, err
	}

	for _, propertyModel := range propertyModels {
		properties = append(properties, propertyModel.fromModel())
	}

	return properties, nil
}

func (r *PropertyRepository) Create(property twomes.Property) (twomes.Property, error) {
	propertyModel := MakePropertyModel(property)
	err := r.db.Create(&propertyModel).Error
	return propertyModel.fromModel(), err
}

func (r *PropertyRepository) Delete(property twomes.Property) error {
	propertyModel := MakePropertyModel(property)
	return r.db.Delete(&propertyModel).Error
}
