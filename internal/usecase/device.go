package usecase

import (
	"context"

	"github.com/dev-shimada/read-write-splitting-playground/internal/domain"
)

type DeviceUsecase struct {
	dr domain.DeviceRepository
}
type DeviceAddInput struct {
	Name   domain.DeviceName
	Status domain.DeviceStatus
}
type DeviceAddOutput struct {
	ID int
}

func NewDeviceUsecase(dr domain.DeviceRepository) DeviceUsecase {
	return DeviceUsecase{
		dr: dr,
	}
}
func (u DeviceUsecase) Add(input DeviceAddInput) (DeviceAddOutput, error) {
	device := domain.NewDevice(
		input.Name,
		input.Status,
	)
	ctx := context.Background()
	id, err := u.dr.Create(ctx, device)
	if err != nil {
		return DeviceAddOutput{}, err
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
