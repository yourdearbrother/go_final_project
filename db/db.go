package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() {
	dbFile := os.Getenv("DB_FILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	log.Printf("Используемый файл БД: %s", dbFile)

	var install bool
}