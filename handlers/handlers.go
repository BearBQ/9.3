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
		SendError(w, http.StatusUnsupportedMediaType, "Expected json")
		return
	}

	var data map[string]string

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		SendError(w, http.StatusBadRequest, "Invalid json")
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
		SendError(w, http.StatusUnprocessableEntity, "validation failed"+err.Error())
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
		SendError(w, http.StatusBadRequest, "ID not found")
		return
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		SendError(w, http.StatusBadRequest, "invalid ID format")
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

	SendError(w, http.StatusNotFound, "Task not found")

}

// Hello возвращает информацию об API
// @Summary Информация об API
// @Description Возвращает описание доступных endpoints
// @Tags info
// @Produce  plain
// @Success 200 {string} string "Описание API"
// @Router / [get]
func (h *MyHandler) Hello(w http.ResponseWriter, r *http.Request) {
	message := "GET /tasks → возвращает список задач.\nPOST /tasks → добавляет задачу.\nDELETE /tasks/{id} → удаляет задачу.\nPUT /tasks/{id} → изменяет задачу"
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

// SendError отправляет ошибку клиенту
func SendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Code:    status,
		Message: message,
	})
}

// PutTask обновляет задачу по ID
// @Summary Обновить задачу
// @Description Обновляет существующую задачу по указанному ID
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param   id path int true "ID задачи"
// @Param   task body object true "Данные для обновления" schema={"type":"object","properties":{"title":{"type":"string","example":"Обновленная задача"},"done":{"type":"boolean","example":true}}}
// @Success 200 {object} models.Task
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 415 {object} models.ErrorResponse
// @Failure 422 {object} models.ErrorResponse
// @Router /tasks/{id} [put]
func (h *MyHandler) PutTask(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != "application/json" {
		SendError(w, http.StatusUnsupportedMediaType, "Expected json")
		return
	}

	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		SendError(w, http.StatusBadRequest, "Invalid json")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		if vars := mux.Vars(r); vars != nil {
			id = vars["id"]
		}
	}
	if id == "" {
		SendError(w, http.StatusBadRequest, "ID not found in URL")
		return
	}

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		SendError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if data["title"] == "" {
		SendError(w, http.StatusBadRequest, "Title is required")
		return
	}

	if data["done"] == "" {
		SendError(w, http.StatusBadRequest, "Done field is required")
		return
	}

	var boolData bool
	boolData, err = strconv.ParseBool(data["done"])
	if err != nil {
		SendError(w, http.StatusBadRequest, "Invalid done value, expected true/false")
		return
	}

	newTask := models.Task{
		ID:    uint(idUint),
		Title: data["title"],
		Done:  boolData,
	}

	validate := validator.New()
	if err := validate.Struct(newTask); err != nil {
		SendError(w, http.StatusUnprocessableEntity, "Validation failed: "+err.Error())
		return
	}

	found := false
	for i := range *h.tasks {
		if (*h.tasks)[i].ID == uint(idUint) {
			(*h.tasks)[i] = newTask
			found = true

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(newTask)
			return
		}
	}

	if !found {
		SendError(w, http.StatusNotFound, "Task not found")
	}
}
