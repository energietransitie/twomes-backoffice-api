package datasourcetype

// An datasourcetype can be a device, cloudfeed or energyquery
type DataSourceType struct {
	ID                    uint             `json:"id"`
	TypeInstanceID        uint             `json:"type_instance_id"`
	Category              Category         `json:"category"`
	Order                 uint             `json:"order"`
	InstallationManualURL string           `json:"installation_url"`
	FAQURL                string           `json:"faq_url"`
	InfoURL               string           `json:"info_url"`
	Precedes              []DataSourceType `json:"precedes"`
	UploadSchedule        string           `json:"upload_schedule"`
	MeasurementSchedule   string           `json:"measurement_schedule"`
	NotificationThreshold string           `json:"notification_threshold"`
}

type Category string

const (
	DeviceType      Category = "device_type"
	CloudFeedType   Category = "cloud_feed_type"
	EnergyQueryType Category = "energy_query_type"
)

func MakeDataSourceType(
	typeInstanceID uint,
	category Category,
	installationManualURL string,
	faqURL string,
	infoURL string,
	precedes []DataSourceType,
	uploadSchedule string,
	measurementSchedule string,
	notificationThreshold string,
) DataSourceType {
	return DataSourceType{
		TypeInstanceID:        typeInstanceID,
		Category:              category,
		InstallationManualURL: installationManualURL,
		FAQURL:                faqURL,
		InfoURL:               infoURL,
		Precedes:              precedes,
		UploadSchedule:        uploadSchedule,
		MeasurementSchedule:   measurementSchedule,
		NotificationThreshold: notificationThreshold,
	}
}
