package apikey

// An APIKey stores api keys for use in the app
type APIKey struct {
	ID      uint   `json:"id"`
	APIName string `json:"api_name"`
	APIKey  string `json:"api_key"`
}

// Create a new APIKey.
func MakeAPIKey(apiName string, apiKey string) APIKey {
	return APIKey{
		APIName: apiName,
		APIKey:  apiKey,
	}
}
