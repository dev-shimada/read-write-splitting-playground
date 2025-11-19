package usecase

import (
	"context"

	"github.com/dev-shimada/read-write-splitting-playground/internal/domain"
	"github.com/dev-shimada/read-write-splitting-playground/internal/infrastructure/repository"
)

type DeviceInput struct {
	Name   domain.DeviceName
	Status domain.DeviceStatus
}
type DeviceOutput struct{}

type DeviceUsecase struct {
	dr repository.DeviceRepository
}

func NewDeviceUsecase(dr repository.DeviceRepository) DeviceUsecase {
	return DeviceUsecase{
		dr: dr,
	}
}
func (u DeviceUsecase) Add(input DeviceInput) (DeviceOutput, error) {
	device := domain.NewDevice(
		input.Name,
		input.Status,
	)
	ctx := context.Background()
	id, err := u.dr.Create(ctx, device)
	if err != nil {
		return DeviceOutput{}, err
	}
	if _, err := u.dr.FindByID(ctx, id); err != nil {
		return DeviceOutput{}, err
	}
	return DeviceOutput{}, nil
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
