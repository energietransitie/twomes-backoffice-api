package energyqueryupload

// An EnergyQueryUploadRepository can load, store and delete uploads.
type EnergyQueryUploadRepository interface {
	Find(energyQueryUpload EnergyQueryUpload) (EnergyQueryUpload, error)
	GetAll() ([]EnergyQueryUpload, error)
	Create(EnergyQueryUpload) (EnergyQueryUpload, error)
	Delete(EnergyQueryUpload) error
}
