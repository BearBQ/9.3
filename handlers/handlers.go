package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"9.3/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type MyHandler struct {
	tasks *[]models.Task
}

func NewMyHandler(tasks *[]models.Task) *MyHandler {
	return &MyHandler{tasks: tasks}
}

// PostTaskFunc создает новую задачу
// @Summary Создать задачу
// @Description Добавляет новую задачу в список
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param   task body object true "Данные задачи" schema={"type":"object","properties":{"title":{"type":"string","example":"Новая задача"}}}
// @Success 200 {object} models.Task
// @Failure 400 {object} models.ErrorResponse
// @Failure 415 {object} models.ErrorResponse
// @Router /tasks [post]
func (h *MyHandler) PostTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusUnsupportedMediaType,
			Message: "Expected json",
		})
		return
	}

	var data map[string]string

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid json",
		})
		return
	}
	id := len(*h.tasks) + 1

	newTask := models.Task{

		ID:    uint(id),
		Title: data["title"],
		Done:  false,
	}

	//Добавляю валидацию структуры
	validate := validator.New()
	if err := validate.Struct(newTask); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: "Validation failed" + err.Error(),
		})
		return
	}

	*h.tasks = append(*h.tasks, newTask)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

// GetTaskFunc возвращает список всех задач
// @Summary Получить все задачи
// @Description Возвращает список всех задач
// @Tags tasks
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Task
// @Router /tasks [get]
func (h *MyHandler) GetTaskFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(h.tasks)
}

// DeleteTaskFunc удаляет задачу по ID
// @Summary Удалить задачу
// @Description Удаляет задачу по указанному ID
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param   id path int true "ID задачи"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /tasks/{id} [delete]
func (h *MyHandler) DeleteTaskFunc(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		//значение для гориллыz
		if vars := mux.Vars(r); vars != nil {
			id = vars["id"]
		}
	}
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "ID not found",
		})
		return
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID format",
		})
		return
	}
	for i, val := range *h.tasks {
		if val.ID == (uint(idUint)) {
			*h.tasks = append((*h.tasks)[:i], (*h.tasks)[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(models.SuccessResponse{
				Message: fmt.Sprintf("Task with ID %d was deleted", val.ID),
			})
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Code:    http.StatusNotFound,
		Message: "Task not found",
	})

}

// Hello возвращает информацию об API
// @Summary Информация об API
// @Description Возвращает описание доступных endpoints
// @Tags info
// @Produce  plain
// @Success 200 {string} string "Описание API"
// @Router / [get]
func (h *MyHandler) Hello(w http.ResponseWriter, r *http.Request) {
	message := "GET /tasks → возвращает список задач.\nPOST /tasks → добавляет задачу.\nDELETE /tasks/{id} → удаляет задачу"
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
