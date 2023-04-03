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

// A Device is collects measurements in a subject's building.
type Device struct {
	ID                   uint       `json:"id"`
	Name                 string     `json:"name"`
	DeviceType           DeviceType `json:"device_type"`
	BuildingID           uint       `json:"building_id"`
	ActivationSecret     string     `json:"activation_secret,omitempty"` // This can be removed if a device uses JWT's too.
	ActivationSecretHash string     `json:"-"`                           // This can be removed if a device uses JWT's too.
	ActivatedAt          *time.Time `json:"activated_at"`
	AuthorizationToken   string     `json:"authorization_token,omitempty"`
	Uploads              []Upload   `json:"uploads,omitempty"`
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
