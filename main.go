package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		fmt.Println(err)
	}
}