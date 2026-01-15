package testutils

import (
	"database/sql"
	"fmt"
	"os"
)

func OpenTestDB() *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return db
}

func Seed(db *sql.DB, sqlFile string) {
	bytes, err := os.ReadFile(sqlFile)
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(bytes)); err != nil {
		panic(err)
	}
}
