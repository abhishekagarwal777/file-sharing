package db

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
)

var DB *sql.DB

func InitDB() {
    var err error
    connStr := "user=User1 password=Pass123 dbname=xyz sslmode=disable"
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    if err := DB.Ping(); err != nil {
        log.Fatal("Cannot connect to DB:", err)
    }
    log.Println("Connected to DB!")
}
