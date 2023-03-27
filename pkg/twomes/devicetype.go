package twomes

import (
	"encoding/json"
	"errors"
	"time"
)

// A DeviceType contains information about a group of devices with the same functionality.
type DeviceType struct {
	ID                    uint       `json:"id"`
	Name                  string     `json:"name"`
	InstallationManualURL string     `json:"installation_manual_url"`
	InfoURL               string     `json:"info_url"`
	Properties            []Property `json:"properties,omitempty"`
	UploadInterval        Duration   `json:"upload_interval"`
}

// Create a new DeviceType.
func MakeDeviceType(name, installationManualURL, infoURL string, properties []Property, uploadInterval Duration) DeviceType {
	return DeviceType{
		Name:                  name,
		InstallationManualURL: installationManualURL,
		InfoURL:               infoURL,
		Properties:            properties,
		UploadInterval:        uploadInterval,
	}
}

var (
	ErrInvalidDuration = errors.New("invalid duration")
)

// Duration is a wrapper around [time.Duration] that can be marshalled and unmarshalled to JSON.
type Duration struct {
	time.Duration
}

// Create a new Duration from a time.Duration.
func MakeDuration(duration time.Duration) Duration {
	return Duration{
		Duration: duration,
	}
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v any
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	value, ok := v.(string)
	if !ok {
		return ErrInvalidDuration
	}

	d.Duration, err = time.ParseDuration(value)
	if err != nil {
		return err
	}

	return nil
}
