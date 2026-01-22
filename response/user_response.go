package response

import "github.com/google/uuid"

type UpdateUser struct {
	Nom    string `json:"nom"`
	Prenom string `json:"prenom"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResume struct {
	ID               uuid.UUID `json:"ID"`
	Nom              string    `json:"nom"`
	TotalTasks       int       `json:"total_tasks"`
	CompletedTasks   int       `json:"compelted_tasks"`
	CompletedPercent float64   `json:"completed_percent"`
}

type UserStat struct {
	TotalUser      int64   `json:"total_user"`
	TotalTasks     int64   `json:"total_tasks"`
	TotalCompleted int64   `json:"total_completed"`
	Rate           float64 `json:"rate"`
}
