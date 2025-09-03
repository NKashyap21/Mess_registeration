package config

import (
	"context"
	"log/slog"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	var err error
	DB, err = pgxpool.New(context.Background(), os.Getenv("DB_URL"))

	if err != nil {
		slog.Error("Failed to initialize DB")
		os.Exit(1)
	}
	// driver, err := pgx.WithInstance(stdlib.OpenDBFromPool(DB), &pgx.Config{})

	// if err != nil {
	// 	slog.Error("Failed to initialize Migrator Driver")
	// 	os.Exit(1)
	// }

	// wd, err := os.Getwd()
	// m, err := migrate.NewWithDatabaseInstance(filepath.Join("file:///", wd, "migrations"), "postgres", driver)
	// if err != nil {
	// 	slog.Error("Failed to initialize New Migrator")
	// 	slog.Error(err.Error())
	// 	os.Exit(1)
	// }
	// if err = m.Up(); err != nil {
	// 	slog.Error("Failed to migrate UP")
	// 	slog.Error(err.Error())
	// }

}
