package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrDeviceNotFound = errors.New("device not found")
)

type DeviceRepository interface {
	Create(ctx context.Context, device Device) (int, error)
	FindByID(ctx context.Context, id int) (Device, error)
}

// Device represents a Device entity
type Device struct {
	id        int
	name      DeviceName
	status    DeviceStatus
	createdAt time.Time
	updatedAt time.Time
}

// NewDevice creates a new device with validation using value objects
func NewDevice(name DeviceName, status DeviceStatus) Device {
	return Device{
		id:        0,
		name:      name,
		status:    status,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
}
func (d Device) ID() int {
	return d.id
}
func (d Device) Name() DeviceName {
	return d.name
}
func (d Device) Status() DeviceStatus {
	return d.status
}
func (d Device) Activate() Device {
	d.status = d.status.Activate()
	return d
}
func (d Device) Deactivate() Device {
	d.status = d.status.Deactivate()
	return d
}
func (d Device) UpdateName(name string) (Device, error) {
	deviceName, err := NewDeviceName(name)
	if err != nil {
		return Device{}, err
	}
	d.name = deviceName
	return d, nil
}
func (d Device) IsActive() bool {
	return d.status.IsActive()
}
