package main

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/dev-shimada/read-write-splitting-playground/internal/infrastructure/database"
	"github.com/dev-shimada/read-write-splitting-playground/internal/infrastructure/repository"
	"github.com/dev-shimada/read-write-splitting-playground/internal/presentation"
	"github.com/dev-shimada/read-write-splitting-playground/internal/usecase"
)

func main() {
	// Database connection settings
	writerDSN := getEnv("WRITER_DSN", "root:password@tcp(writer-db:3306)/mydb?parseTime=true")
	readerDSN := getEnv("READER_DSN", "root:password@tcp(reader-db:3306)/mydb?parseTime=true")

	// Initialize writer DB connection
	writerDB, err := sqlx.Connect("mysql", writerDSN)
	if err != nil {
		log.Fatalf("Failed to connect to writer DB: %v", err)
	}
	defer writerDB.Close()

	// Initialize reader DB connection
	readerDB, err := sqlx.Connect("mysql", readerDSN)
	if err != nil {
		log.Fatalf("Failed to connect to reader DB: %v", err)
	}
	defer readerDB.Close()

	// Initialize DBAccessor
	dbAccessor := database.NewDBAccessor(writerDB, readerDB)

	// Initialize repository layer
	deviceRepo := repository.NewDeviceRepository(dbAccessor)

	// Initialize usecase layer
	deviceUsecase := usecase.NewDeviceUsecase(deviceRepo)

	// Initialize CLI presentation layer
	cli := presentation.NewCLI(deviceUsecase)

	// Run CLI application
	exitCode := cli.Run(os.Args)
	os.Exit(exitCode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
