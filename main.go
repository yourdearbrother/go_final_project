package main

import (
	"fmt"
	"net/http"
	"go_final_project/db"
)

func main() {

	db.Init()
	defer db.DB.Close()

	http.Handle("/", http.FileServer(http.Dir("./web")))
	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		fmt.Println(err)
	}
}