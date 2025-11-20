package presentation

import (
	"fmt"

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
	addOutput, err := c.du.Add(usecase.DeviceAddInput{
		Name:   name,
		Status: status,
	})
	if err != nil {
		return err
	}
	fmt.Printf("addOutput: %+v\n", addOutput)
	findOutput, err := c.du.Find(usecase.DeviceFindInput(addOutput))
	if err != nil {
		return err
	}
	fmt.Printf("findOutput: %+v\n", findOutput)
	return nil
}
