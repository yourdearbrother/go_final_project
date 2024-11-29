package main

import (
	"log"
	"net/http"

	"go_final_project/db"
	"go_final_project/handlers"
)

func main() {

	db.Init()
	defer db.DB.Close()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("web/")))
	mux.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	mux.HandleFunc("/api/task", handlers.TaskHandler)
	mux.HandleFunc("/api/tasks", handlers.GetTasksHandler)
	mux.HandleFunc("/api/task/done", handlers.HandleCompleteTask)

	err := http.ListenAndServe(":7540", mux)
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return
	}
}
