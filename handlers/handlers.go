package handlers

import (
	"encoding/json"
	"net/http"

	"9.3/models"
)

type MyHandler struct {
	tasks []models.Task
}

func NewMyHandler(tasks []models.Task) *MyHandler {
	return &MyHandler{tasks: tasks}
}

func (h *MyHandler) PostTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != "application/json" {
		http.Error(w, "Expected json", http.StatusUnsupportedMediaType)
		return
	}

	var data models.IncomingTask
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
	}
	newTask := models.Task{
		ID:    data.ID,
		Title: data.Task,
		Done:  false,
	}
	h.tasks = append(h.tasks, newTask)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

func (h *MyHandler) GetTaskFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(h.tasks)
}
