package models

// Task представляет модель задачи
type Task struct {
	ID    uint   `json:"id" validate:"required,min=1"`
	Title string `json:"title" validate:"required,min=3,max=255"`
	Done  bool   `json:"done"`
}

// ErrorResponse представляет модель ошибки
type ErrorResponse struct {
	// Код ошибки
	Code int `json:"code" example:"400"`
	// Сообщение об ошибке
	Message string `json:"message" example:"Invalid input"`
}
type SuccessResponse struct {
	// Сообщение об успехе
	Message string `json:"message" example:"Task was deleted successfully"`
}
