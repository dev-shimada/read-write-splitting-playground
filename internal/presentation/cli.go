package presentation

import (
	"github.com/dev-shimada/read-write-splitting-playground/internal/domain"
	"github.com/dev-shimada/read-write-splitting-playground/internal/usecase"
)

type CLI struct {
	du usecase.DeviceUsecase
}

func NewCLI(du usecase.DeviceUsecase) CLI {
	return CLI{
		du: du,
	}
}

func (c CLI) Run() error {
	name, err := domain.NewDeviceName("Device001")
	if err != nil {
		return err
	}
	status, err := domain.NewDeviceStatus("active")
	if err != nil {
		return err
	}
	_, err = c.du.Add(usecase.DeviceInput{
		Name:   name,
		Status: status,
	})
	if err != nil {
		return err
	}
	return nil
}
