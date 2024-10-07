package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func Init() {
	dbFile := os.Getenv("DB_FILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	log.Printf("Используемый файл БД: %s", dbFile)

	var install bool

	_, err := os.Stat(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			log.Fatalf("Не удалось проверить существует ли файл БД: %v", err)
		}
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	DB = db

	if install {
		log.Println("Создаем БД")


		_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL CHECK(length(date) = 8),
            title TEXT NOT NULL,
            comment TEXT,
            repeat TEXT CHECK(length(repeat) <= 128)
        );`)
		if err != nil {
			log.Fatalf("Ошибка при создании таблицы: %v", err)
		}

		_, err = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`)
		if err != nil {
			log.Fatalf("Ошибка при создании индекса: %v", err)
		}

		log.Println("БД создана")
	} else {
		log.Println("БД уже существует")
	}
}