// NOTE:
// This DB package is optional and NOT part of the default datasource runtime.
// Persistence is handled downstream by Ingestor / API Gateway.
// This package is retained for local debugging or future extensions.package db

import (
    "database/sql"
    "fmt"
    "os"

    _ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("POSTGRES_HOST"),
        os.Getenv("POSTGRES_PORT"),
        os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_PASSWORD"),
        os.Getenv("POSTGRES_DB"),
    )

    return sql.Open("postgres", dsn)
}
