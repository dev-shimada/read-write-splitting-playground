package domain

import (
	"errors"
	"time"
)

var (
	ErrDeviceNotFound = errors.New("device not found")
)

// Device is an alias for the unexported device struct
type Device = device

// Device represents a device entity
type device struct {
	id        int64
	name      deviceName
	status    deviceStatus
	createdAt time.Time
	updatedAt time.Time
}

// NewDevice creates a new device with validation using value objects
func NewDevice(name string, status string) (*device, error) {
	deviceName, err := NewDeviceName(name)
	if err != nil {
		return nil, err
	}

	deviceStatus, err := NewDeviceStatus(status)
	if err != nil {
		return nil, err
	}

	return &device{
		id:        0,
		name:      deviceName,
		status:    deviceStatus,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

// Validate validates the device fields
func (d *device) Validate() error {
	// Validation is already done in value objects during construction
	return nil
}

// Activate sets the device status to active using value object
func (d *device) Activate() error {
	d.status = d.status.Activate()
	return nil
}

// Deactivate sets the device status to inactive using value object
func (d *device) Deactivate() error {
	d.status = d.status.Deactivate()
	return nil
}

// UpdateName updates the device name with validation using value object
func (d *device) UpdateName(name string) error {
	deviceName, err := NewDeviceName(name)
	if err != nil {
		return err
	}
	d.name = deviceName
	return nil
}

// IsActive returns true if the device is active
func (d *device) IsActive() bool {
	return d.status.IsActive()
}
