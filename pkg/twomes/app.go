package twomes

// An App can provision a [Device] in a [Building].
type App struct {
	ID                      uint   `json:"id"`
	Name                    string `json:"name"`
	ProvisioningURLTemplate string `json:"provisioning_url_template"`
}

// Create a new app.
func MakeApp(name, provisioningURLTemplate string) App {
	return App{
		Name:                    name,
		ProvisioningURLTemplate: provisioningURLTemplate,
	}
}
