package main

import (
	"log/slog"

	_ "github.com/go-sql-driver/mysql"

	"github.com/dev-shimada/read-write-splitting-playground/internal/infrastructure/database"
	"github.com/dev-shimada/read-write-splitting-playground/internal/infrastructure/repository"
	"github.com/dev-shimada/read-write-splitting-playground/internal/presentation"
	"github.com/dev-shimada/read-write-splitting-playground/internal/usecase"
)

func main() {
	db := database.NewDBAccessor()
	if err := database.Migrate(db); err != nil {
		slog.Error("failed to migrate database", "error", err)
		panic(err)
	}
	dr := repository.NewDeviceRepository(db)
	du := usecase.NewDeviceUsecase(dr, db)
	cli := presentation.NewCLI(du)
	if err := cli.Run(); err != nil {
		slog.Error("failed to run CLI", "error", err)
		panic(err)
	}
}
