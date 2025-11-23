package usecase

import (
	"context"
	"fmt"

	"github.com/dev-shimada/read-write-splitting-playground/internal/domain"
)

type transactioner interface {
	Transaction(context.Context, func(context.Context) error) error
}

type DeviceUsecase struct {
	dr  domain.DeviceRepository
	dba transactioner
}
type DeviceAddInput struct {
	Name   domain.DeviceName
	Status domain.DeviceStatus
}
type DeviceAddOutput struct {
	ID int
}

func NewDeviceUsecase(dr domain.DeviceRepository, dba transactioner) DeviceUsecase {
	return DeviceUsecase{
		dr:  dr,
		dba: dba,
	}
}
func (u DeviceUsecase) Add(input DeviceAddInput) (DeviceAddOutput, error) {
	device := domain.NewDevice(
		input.Name,
		input.Status,
	)
	ctx := context.Background()

	var id int
	if txErr := u.dba.Transaction(ctx, func(txCtx context.Context) error {
		var err error
		id, err = u.dr.Create(txCtx, device)
		if err != nil {
			return fmt.Errorf("failed to create device: %w", err)
		}
		return nil
	}); txErr != nil {
		return DeviceAddOutput{}, fmt.Errorf("failed to add device in transaction: %w", txErr)
	}

	return DeviceAddOutput{ID: id}, nil
}

type DeviceFindInput struct {
	ID int
}
type DeviceFindOutput struct {
	Device domain.Device
}

func (u DeviceUsecase) Find(input DeviceFindInput) (DeviceFindOutput, error) {
	ctx := context.Background()
	device, err := u.dr.FindByID(ctx, input.ID)
	if err != nil {
		return DeviceFindOutput{}, err
	}
	return DeviceFindOutput{Device: device}, nil
}

// type transactioner interface {
// 	Transaction(context.Context, func(context.Context) error) error
// }

// type unreadCountRepository interface {
// 	GetForUpdate(context.Context) (int, error)
// 	Increment(context.Context) error
// }

// type notifyUserDetailsUsecase struct {
// 	dba                   transactioner
// 	unreadCountRepository unreadCountRepository
// }

// func (u *notifyUserDetailsUsecase) NotifyUserDetails(
// 	ctx context.Context,
// ) error {
// 	// txCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	// defer cancel()

// 	if txErr := u.dba.Transaction(ctx, func(txCtx context.Context) error {
// 		if _, err := u.unreadCountRepository.GetForUpdate(txCtx); err != nil {
// 			return fmt.Errorf("failed to lock unread count: %w", err)
// 		}

// 		if err := u.unreadCountRepository.Increment(txCtx); err != nil {
// 			return fmt.Errorf("failed to increment unread count: %w", err)
// 		}

// 		return nil
// 	}); txErr != nil {
// 		return fmt.Errorf("failed to dba.Transaction: %w", txErr)
// 	}

// 	return nil
// }
