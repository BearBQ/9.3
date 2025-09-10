package main

import (
	"log"
	"net/http"
	"time"

	"9.3/handlers"
	"9.3/models"
	"github.com/gorilla/mux"
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
	go func() {
		mux := http.NewServeMux()
		handler := handlers.NewMyHandler(taskDesk)
		mux.HandleFunc("GET /", handler.Hello)
		mux.HandleFunc("GET /tasks", handler.GetTaskFunc)
		mux.HandleFunc("POST /tasks", handler.PostTaskFunc)
		mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTaskFunc)
		mux.HandleFunc("GET /swagger/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "http://localhost:8081/swagger/", http.StatusMovedPermanently)
		})
		log.Println("server http in go routine is started")
		err := http.ListenAndServe(":8080", mux)
		if err != nil {
			log.Fatalf("http server error: %v", err)
		}
	}()

	//реализую через гориллу
	r := mux.NewRouter()
	srv := &http.Server{
		Addr: "0.0.0.0:8081",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.

	}
	handler := handlers.NewMyHandler(taskDesk)

	// API routes
	r.HandleFunc("/", handler.Hello).Methods("GET")
	r.HandleFunc("/tasks", handler.GetTaskFunc).Methods("GET")
	r.HandleFunc("/tasks", handler.PostTaskFunc).Methods("POST")
	r.HandleFunc("/tasks/{id}", handler.DeleteTaskFunc).Methods("DELETE")

	// Swagger UI route
	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs"))))

	log.Println("server gorilla is started")
	log.Println("Swagger UI available at: http://localhost:8081/swagger/index.html")
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
