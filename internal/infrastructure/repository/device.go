package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type dbAccessor interface {
	Transaction(ctx context.Context, txFunc func(context.Context) error) error
	Exec(ctx context.Context, query string, arg any) (sql.Result, error)
	Query(ctx context.Context, query string, arg any) (*sqlx.Rows, error)
}

type deviceRepository struct {
	db dbAccessor
}

// NewDeviceRepository creates a new device repository
func NewDeviceRepository(db dbAccessor) deviceRepository {
	return deviceRepository{
		db: db,
	}
}
