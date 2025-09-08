package models

type Task struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type IncomingTask struct {
	ID   uint   `json:"id"`
	Task string `json:"task"`
}
