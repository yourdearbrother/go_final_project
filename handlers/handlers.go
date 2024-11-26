package handlers

import (
	"net/http"
	"go_final_project/utils"
	"go_final_project/db"
	"time"
	"encoding/json"
	"strconv"
	"fmt"
)

func isDate(str string) bool {
	_, err := time.Parse("02.01.2006", str)
	return err == nil
}

func convertToDate(str string) string {
	date, _ := time.Parse("02.01.2006", str)
	return date.Format("20060102")
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	const layout = "20060102"
	now, err := time.Parse(layout, nowStr)
	if err != nil {
		http.Error(w, "время не удалось преобразовать", http.StatusBadRequest)
		return
	}

	nextDate, err := utils.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}

type Task struct {
	ID      string `db:"id" json:"id"`
	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat"`
}

// Переключаем методы
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleCreateTask(w, r)
	case http.MethodPut:
		handleUpdateTask(w, r)
	case http.MethodGet:
		handleGetTask(w, r)
	case http.MethodDelete:
		handleDeleteTask(w, r)
	default:
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
	}
}

func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response := map[string]string{"error": "ошибка десериализации JSON"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if task.Title == "" {
		response := map[string]string{"error": "не указан заголовок задачи"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	const layout = "20060102"
	now := time.Now()
	nowStr := now.Format(layout)

	if task.Date == "" {
		task.Date = nowStr
	} else {
		parsedDate, err := time.Parse(layout, task.Date)
		if err != nil {
			response := map[string]string{"error": "дата указана в неверном формате"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if task.Date == nowStr {
		} else if parsedDate.Before(now) {
			if task.Repeat == "" {
				task.Date = nowStr
			} else {
				nextDate, err := utils.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					response := map[string]string{"error": err.Error()}
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(response)
					return
				}
				task.Date = nextDate
			}
		}
	}

	id, err := db.AddTask(task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		response := map[string]string{"error": err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]string{"id": fmt.Sprintf("%d", id)}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response := map[string]string{"error": "ошибка десериализации JSON"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if task.ID == "" {
		response := map[string]string{"error": "не указан идентификатор задачи"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if task.Title == "" {
		response := map[string]string{"error": "не указан заголовок задачи"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	const layout = "20060102"
	now := time.Now()
	nowStr := now.Format(layout)

	if task.Date == "" {
		task.Date = nowStr
	} else {
		parsedDate, err := time.Parse(layout, task.Date)
		if err != nil {
			response := map[string]string{"error": "дата указана в неверном формате"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if task.Date == nowStr {
		} else if parsedDate.Before(now) {
			if task.Repeat == "" {
				task.Date = nowStr
			} else {
				nextDate, err := utils.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					response := map[string]string{"error": err.Error()}
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(response)
					return
				}
				task.Date = nextDate
			}
		}
	}

	err := db.UpdateTask(task.ID, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		response := map[string]string{"error": err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, `{"error":"не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"неправильный идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := db.GetTaskByID(id)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(task)
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, `{"error":"не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"неправильный идентификатор"}`, http.StatusBadRequest)
		return
	}

	err = db.DeleteTask(id)
	if err != nil {
		http.Error(w, `{"error":"не удалось удалить задачу"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{})
}


func HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, `{"error":"не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"некорректный идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := db.GetTaskByID(id)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			http.Error(w, `{"error":"не удалось удалить задачу"}`, http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := utils.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"не удалось получить следующую дату"}`, http.StatusBadRequest)
			return
		}

		err = db.UpdateTask(task.ID, nextDate, task.Title, task.Comment, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"не удалось обновить задачу"}`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{})
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {

	searchParam := r.URL.Query().Get("search")

	var tasks []db.Task
	var err error

	if searchParam != "" {
		if isDate(searchParam) {
			tasks, err = db.GetTasksByDate(convertToDate(searchParam))
		} else {
			tasks, err = db.SearchTasks(searchParam)
		}
	} else {
		tasks, err = db.GetTasks()
	}

	if err != nil {
		http.Error(w, `{"error":"не удалось полученить список задач"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"tasks": tasks,
	}

	if tasks == nil {
		response["tasks"] = []db.Task{}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

func GetTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, `{"error":"не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"неправильный идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := db.GetTaskByID(id)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(task)
}