package main

import (
	"log"
	"net/http"

	"9.3/handlers"
	"9.3/models"
)

func main() {
	taskDesk := []models.Task{
		{ID: 1,
			Title: "first task",
			Done:  true},
		{ID: 2,
			Title: "second task",
			Done:  true},
		{ID: 3,
			Title: "third task",
			Done:  false},
		{ID: 4,
			Title: "fourth task",
			Done:  true},
		{ID: 5,
			Title: "fifth task",
			Done:  true},
		{ID: 6,
			Title: "sixth task",
			Done:  false},
	}
	//версия на стандартном http
	mux := http.NewServeMux()
	handler := handlers.NewMyHandler(taskDesk)
	mux.HandleFunc("GET /", handler.Hello)
	mux.HandleFunc("GET /tasks", handler.GetTaskFunc)
	mux.HandleFunc("POST /tasks", handler.PostTaskFunc)
	mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTaskFunc)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("http server error: %v", err)
	}

}
