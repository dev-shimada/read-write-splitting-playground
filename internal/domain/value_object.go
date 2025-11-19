package domain

import (
	"errors"
	"fmt"
	"slices"
)

var (
	ErrEmptyDeviceName    = errors.New("device name cannot be empty")
	ErrDeviceNameTooLong  = errors.New("device name must be less than 256 characters")
	ErrInvalidStatusValue = errors.New("invalid device status value")
)

type DeviceName struct {
	value string
}

func NewDeviceName(name string) (DeviceName, error) {
	if name == "" {
		return DeviceName{}, ErrEmptyDeviceName
	}
	if len(name) > 255 {
		return DeviceName{}, ErrDeviceNameTooLong
	}
	return DeviceName{value: name}, nil
}

func (d DeviceName) Value() string {
	return d.value
}
func (d DeviceName) String() string {
	return d.value
}
func (d DeviceName) Equals(other DeviceName) bool {
	return d.value == other.value
}

type DeviceStatus struct {
	value string
}

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

func NewDeviceStatus(status string) (DeviceStatus, error) {
	if !slices.Contains([]string{StatusActive, StatusInactive}, status) {
		return DeviceStatus{}, fmt.Errorf("%w: %s", ErrInvalidStatusValue, status)
	}
	return DeviceStatus{value: status}, nil
}

func (d DeviceStatus) Value() string {
	return d.value
}
func (d DeviceStatus) String() string {
	return d.value
}
func (d DeviceStatus) Equals(other DeviceStatus) bool {
	return d.value == other.value
}
func (d DeviceStatus) IsActive() bool {
	return d.value == StatusActive
}
func (d DeviceStatus) IsInactive() bool {
	return d.value == StatusInactive
}
func (d DeviceStatus) Activate() DeviceStatus {
	return DeviceStatus{value: StatusActive}
}
func (d DeviceStatus) Deactivate() DeviceStatus {
	return DeviceStatus{value: StatusInactive}
}
