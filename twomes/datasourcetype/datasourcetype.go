package datasourcetype

// An datasourcetype can be a device, cloudfeed or energyquery
type DataSourceType struct {
	ID                    uint             `json:"id"`
	TypeSourceID          uint             `json:"type_source_id"`
	Type                  CategoryType     `json:"type"`
	InstallationManualURL string           `json:"installation_manual_url"`
	InfoURL               string           `json:"info_url"`
	Precedes              []DataSourceType `json:"precedes"`
	UploadSchedule        []string         `json:"upload_schedule"`
	MeasurementSchedule   []string         `json:"measurement_schedule"`
	NotificationThreshold string           `json:"notification_threshold"`
}

type CategoryType string

const (
	Device_Type       CategoryType = "device_type"
	Cloud_Feed_Type   CategoryType = "cloud_feed_type"
	Energy_Query_Type CategoryType = "energy_query_type"
)

func MakeDataSourceType(
	dataSourceTypeID uint,
	itemType CategoryType,
	installationManualURL string,
	infoURL string,
	precedes []DataSourceType,
	uploadSchedule []string,
	measurementSchedule []string,
	notificationThreshold string,
) DataSourceType {
	return DataSourceType{
		TypeSourceID:          dataSourceTypeID,
		Type:                  itemType,
		InstallationManualURL: installationManualURL,
		InfoURL:               infoURL,
		Precedes:              precedes,
		UploadSchedule:        uploadSchedule,
		MeasurementSchedule:   measurementSchedule,
		NotificationThreshold: notificationThreshold,
	}
}
