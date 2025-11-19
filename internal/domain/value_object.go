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

// DeviceName is a value object representing a device name
type deviceName struct {
	value string
}

// NewDeviceName creates a new DeviceName value object with validation
func NewDeviceName(name string) (deviceName, error) {
	if name == "" {
		return deviceName{}, ErrEmptyDeviceName
	}

	if len(name) > 255 {
		return deviceName{}, ErrDeviceNameTooLong
	}

	return deviceName{value: name}, nil
}

// Value returns the string value of the device name
func (d deviceName) Value() string {
	return d.value
}

// String implements the Stringer interface
func (d deviceName) String() string {
	return d.value
}

// Equals checks if two DeviceName values are equal
func (d deviceName) Equals(other deviceName) bool {
	return d.value == other.value
}

// DeviceStatus is a value object representing a device status
type deviceStatus struct {
	value string
}

// Valid status constants
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

// ValidStatuses returns a list of valid device statuses
var ValidStatuses = []string{StatusActive, StatusInactive}

// NewDeviceStatus creates a new DeviceStatus value object with validation
func NewDeviceStatus(status string) (deviceStatus, error) {
	if !isValidStatus(status) {
		return deviceStatus{}, fmt.Errorf("%w: %s", ErrInvalidStatusValue, status)
	}
	return deviceStatus{value: status}, nil
}

// isValidStatus checks if the status is valid
func isValidStatus(status string) bool {
	return slices.Contains(ValidStatuses, status)
}

// Value returns the string value of the device status
func (d deviceStatus) Value() string {
	return d.value
}

// String implements the Stringer interface
func (d deviceStatus) String() string {
	return d.value
}

// Equals checks if two DeviceStatus values are equal
func (d deviceStatus) Equals(other deviceStatus) bool {
	return d.value == other.value
}

// IsActive returns true if the status is active
func (d deviceStatus) IsActive() bool {
	return d.value == StatusActive
}

// IsInactive returns true if the status is inactive
func (d deviceStatus) IsInactive() bool {
	return d.value == StatusInactive
}

// Activate returns a new DeviceStatus with active status
func (d deviceStatus) Activate() deviceStatus {
	return deviceStatus{value: StatusActive}
}

// Deactivate returns a new DeviceStatus with inactive status
func (d deviceStatus) Deactivate() deviceStatus {
	return deviceStatus{value: StatusInactive}
}
