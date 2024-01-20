package entities

import (
	"context"
	"notification-api/app/infrastructure/mysql"
	"time"
)

type Job struct {
	ID           int       `json:"id"`
	RequestID    string    `json:"request_id"`
	Caller       string    `json:"caller"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Payload      string    `json:"payload"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
	FinishedAt   time.Time `json:"finished_at"`
}

func (Job) TableName() string {
	return "job"
}

func (e *Job) Update(dbRepo *mysql.DBRepository, ctx context.Context) error {
	db := dbRepo.NotificationDB
	return dbRepo.Update(db, ctx, e)
}

func (e *Job) Insert(dbRepo *mysql.DBRepository, ctx context.Context) error {
	db := dbRepo.NotificationDB
	return dbRepo.Insert(db, ctx, e)
}
