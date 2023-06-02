package twomes

// An App can provision a [Device] in a [Building].
type App struct {
	ID                      uint   `json:"id"`
	Name                    string `json:"name"`
	ProvisioningURLTemplate string `json:"provisioning_url_template"`
	OauthRedirectURL        string `json:"oauth_redirect_url"`
}

// Create a new app.
func MakeApp(name, provisioningURLTemplate, oauthRedirectURL string) App {
	return App{
		Name:                    name,
		ProvisioningURLTemplate: provisioningURLTemplate,
		OauthRedirectURL:        oauthRedirectURL,
	}
}
