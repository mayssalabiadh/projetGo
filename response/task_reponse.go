package response

import "projet1/models"

type CompletionRate struct {
	Rate string `json:"completion_rate"`
}

type TaskResult struct {
	Tasks []models.Task
	Err   error
}
