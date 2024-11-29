package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const Layout = "20060102"

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

func AddTask(date, title, comment, repeat string) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, date, title, comment, repeat)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetTasksByDate(date string) ([]Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT 30`
	rows, err := DB.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var id int64
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		task.ID = strconv.Itoa(int(id))
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTasks() ([]Task, error) {
	now := time.Now().Format(Layout)
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date >= ? ORDER BY date`
	rows, err := DB.Query(query, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var id int64
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		task.ID = strconv.Itoa(int(id))
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTaskByID(id int64) (Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := DB.QueryRow(query, id)

	var task Task
	var taskID int64
	err := row.Scan(&taskID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, fmt.Errorf("задача не найдена")
		}
		return Task{}, err
	}
	task.ID = strconv.Itoa(int(id))
	return task, nil
}

func UpdateTask(id, date, title, comment, repeat string) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, date, title, comment, repeat, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func DeleteTask(id int64) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func SearchTasks(search string) ([]Task, error) {
	searchTerm := "%" + search + "%"
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date`
	rows, err := DB.Query(query, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var id int64
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		task.ID = strconv.Itoa(int(id))
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
