package repository

import (
	"context"
	"database/sql"

	"github.com/dev-shimada/read-write-splitting-playground/internal/domain"
	"github.com/jmoiron/sqlx"
)

type dbAccessor interface {
	Transaction(ctx context.Context, txFunc func(context.Context) error) error
	Exec(ctx context.Context, query string, arg any) (sql.Result, error)
	Query(ctx context.Context, query string, arg any) (*sqlx.Rows, error)
}

type DeviceRepository = deviceRepository
type deviceRepository struct {
	db dbAccessor
}

func NewDeviceRepository(db dbAccessor) deviceRepository {
	return deviceRepository{
		db: db,
	}
}
func (r deviceRepository) Create(ctx context.Context, device domain.Device) (int, error) {
	result, err := r.db.Exec(ctx,
		"INSERT INTO devices (name, status) VALUES (:name, :status)",
		struct {
			Name   string `db:"name"`
			Status string `db:"status"`
		}{
			Name:   device.Name().Value(),
			Status: device.Status().Value(),
		},
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
func (r deviceRepository) FindByID(ctx context.Context, id int) (domain.Device, error) {
	row, err := r.db.Query(ctx,
		"SELECT id, name, status FROM devices WHERE id = :id",
		struct {
			ID int `db:"id"`
		}{
			ID: id,
		},
	)
	if err != nil {
		return domain.Device{}, err
	}
	defer row.Close()
	if row.Next() {
		var id int
		var name string
		var status string
		if err := row.Scan(&id, &name, &status); err != nil {
			return domain.Device{}, err
		}
		deviceName, err := domain.NewDeviceName(name)
		if err != nil {
			return domain.Device{}, err
		}
		deviceStatus, err := domain.NewDeviceStatus(status)
		if err != nil {
			return domain.Device{}, err
		}
		return domain.NewDevice(
			deviceName,
			deviceStatus,
		), nil
	}
	return domain.Device{}, sql.ErrNoRows
}
