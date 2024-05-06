package apikey

// An APIKeyRepository can load, store and delete API keys.
type APIKeyRepository interface {
	Find(apiKey APIKey) (APIKey, error)
	GetAll() ([]APIKey, error)
	Create(APIKey) (APIKey, error)
	Delete(APIKey) error
}
