package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"9.3/models"
	"github.com/gorilla/mux"
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

	var data map[string]string

	err := json.NewDecoder(r.Body).Decode(&data)
	title, ok := data["title"]
	if !ok || title == "" || err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}
	id := len(h.tasks) + 1

	newTask := models.Task{

		ID:    uint(id),
		Title: title,
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

func (h *MyHandler) DeleteTaskFunc(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		//значение для гориллыz
		if vars := mux.Vars(r); vars != nil {
			id = vars["id"]
		}
	}
	if id == "" {
		http.Error(w, "ID not found", http.StatusBadRequest)
		return
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
	}
	for i, val := range h.tasks {
		if val.ID == (uint(idUint)) {
			h.tasks = append(h.tasks[:i], h.tasks[i+1:]...)
			message := fmt.Sprintf("Task with ID %d was daleted", val.ID)
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": message,
			})
			return
		}
	}
	http.Error(w, "task not found", http.StatusNotFound)

}

func (h *MyHandler) Hello(w http.ResponseWriter, r *http.Request) {
	message := "GET /tasks → возвращает список задач.\nPOST /tasks → добавляет задачу.\nDELETE /tasks/{id} → удаляет задачу"
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
