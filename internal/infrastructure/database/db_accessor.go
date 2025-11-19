package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type (
	DBAccessor struct {
		writer *sqlx.DB
		reader *sqlx.DB
	}

	abstractSqlxDB interface {
		BindNamed(query string, arg any) (q string, args []any, err error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	}

	ctxKeyTx struct{}
)

func NewDBAccessor() *DBAccessor {
	writer, err := sqlx.Connect("sqlite3", "file:device?mode=memory&cache=shared&_fk=1")
	if err != nil {
		panic(err)
	}
	reader, err := sqlx.Connect("sqlite3", "file:device?mode=memory&cache=shared&_fk=1")
	if err != nil {
		panic(err)
	}
	return &DBAccessor{
		writer: writer,
		reader: reader,
	}
}

func (dba *DBAccessor) Transaction(ctx context.Context, txFunc func(context.Context) error) (err error) {
	tx, err := dba.writer.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			if err := tx.Rollback(); err != nil {
				// ログ出力
			}
			err = fmt.Errorf("panicked on execution txFunc: %v", r)
		}
	}()

	txCtx := withTxDB(ctx, tx)

	if txErr := txFunc(txCtx); txErr != nil {
		if err := tx.Rollback(); err != nil {
			// ログ出力
		}

		return txErr
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (dba *DBAccessor) Exec(
	ctx context.Context,
	query string,
	arg any,
) (sql.Result, error) {
	tx, err := getTxFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to getTxFromContext: %w", err)
	}

	var da abstractSqlxDB
	if tx != nil {
		da = tx
	} else {
		da = dba.writer
	}

	// prepare
	namedQuery, namedArgs, err := da.BindNamed(query, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to BindNamed: %w", err)
	}
	namedQuery, namedArgs, err = sqlx.In(namedQuery, namedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to sqlx.In: %w", err)
	}

	result, err := da.ExecContext(ctx, namedQuery, namedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to run exec: %w", err)
	}

	return result, nil
}

func (dba *DBAccessor) Query(
	ctx context.Context,
	query string,
	arg any,
) (*sqlx.Rows, error) {
	tx, err := getTxFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to getTxFromContext: %w", err)
	}

	var da abstractSqlxDB
	if tx != nil {
		da = tx
	} else {
		da = dba.reader
	}

	// prepare
	namedQuery, namedArgs, err := da.BindNamed(query, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to BindNamed: %w", err)
	}
	namedQuery, namedArgs, err = sqlx.In(namedQuery, namedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to sqlx.In: %w", err)
	}

	rows, err := da.QueryxContext(ctx, namedQuery, namedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %w", err)
	}

	return rows, nil
}

func withTxDB(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, &ctxKeyTx{}, tx)
}

func getTxFromContext(ctx context.Context) (*sqlx.Tx, error) {
	if v := ctx.Value(&ctxKeyTx{}); v != nil {
		tx, ok := v.(*sqlx.Tx)
		if !ok {
			return nil, fmt.Errorf("failed to convert to *sqlx.Tx: %v", v)
		}

		return tx, nil
	}

	return nil, nil
}
