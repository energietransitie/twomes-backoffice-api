package upload

// An UploadRepository can load, store and delete uploads.
type UploadRepository interface {
	Find(Upload Upload) (Upload, error)
	GetAll() ([]Upload, error)
	Create(Upload) (Upload, error)
	Delete(Upload) error
	GetLatestUploadForDeviceWithID(id uint) (Upload, error)
}
