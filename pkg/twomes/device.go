package twomes

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDeviceActivationSecretIncorrect = errors.New("device activation_secret is incorrect")
)

// DeviceHealth describes the health of the device.
// There are 4 health states:
//
//	unactivated = device is created, but not activated yet.
//	activated = device is activated, but no uploads have been received.
//	healthy = device is activated and uploads have been received according to schedule.
//	unhealthy = device is activated, but upload have not been received according to schedule.
type DeviceHealth string

const (
	DeviceHealthUnactivated DeviceHealth = "unactivated"
	DeviceHealthActivated   DeviceHealth = "activated"
	DeviceHealthHealthy     DeviceHealth = "healthy"
	DeviceHealthUnhealthy   DeviceHealth = "unhealthy"
)

// A Device is collects measurements in a subject's building.
type Device struct {
	ID                   uint         `json:"id"`
	Name                 string       `json:"name"`
	DeviceType           DeviceType   `json:"device_type"`
	BuildingID           uint         `json:"building_id"`
	ActivationSecret     string       `json:"activation_secret,omitempty"` // This can be removed if a device uses JWT's too.
	ActivationSecretHash string       `json:"-"`                           // This can be removed if a device uses JWT's too.
	ActivatedAt          *time.Time   `json:"activated_at"`
	AuthorizationToken   string       `json:"authorization_token,omitempty"`
	Uploads              []Upload     `json:"uploads,omitempty"`
	Health               DeviceHealth `json:"health,omitempty"`
}

// Create a new Device.
func MakeDevice(name string, deviceType DeviceType, buildingID uint, activationSecret string) Device {
	activationSecretHash, err := bcrypt.GenerateFromPassword([]byte(activationSecret), 12)
	if err != nil {
		logrus.Error("a device was created, but activationSecretHash could not be generated")
	}

	return Device{
		Name:                 name,
		DeviceType:           deviceType,
		BuildingID:           buildingID,
		ActivationSecretHash: string(activationSecretHash),
	}
}

// Activate a device.
func (d *Device) Activate(activationSecret string) error {
	if activationSecret == "" || bcrypt.CompareHashAndPassword([]byte(d.ActivationSecretHash), []byte(activationSecret)) != nil {
		return ErrDeviceActivationSecretIncorrect
	}

	now := time.Now().UTC()
	d.ActivatedAt = &now

	return nil
}

// Add an upload to a device.
func (d *Device) AddUpload(upload Upload) {
	d.Uploads = append(d.Uploads, upload)
}

func (d *Device) UpdateHealth() {
	if d.ActivatedAt == nil {
		d.Health = DeviceHealthUnactivated
		return
	}

	if len(d.Uploads) <= 0 {
		d.Health = DeviceHealthActivated
		return
	}

	lastUpload := d.Uploads[len(d.Uploads)-1]
	if time.Since(lastUpload.ServerTime) < d.DeviceType.UploadInterval.Duration {
		d.Health = DeviceHealthHealthy
		return
	}

	d.Health = DeviceHealthUnhealthy
}
