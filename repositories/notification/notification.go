package notification

import (
	"context"
	"fmt"
	"go.elastic.co/apm/v2"
	"notification-api/app/infrastructure/mysql"
	"notification-api/domain/entities"
	"time"
)

type Job struct {
	ID           int       `gorm:"column:id;autoIncrement"`
	RequestID    string    `gorm:"column:request_id"`
	Caller       string    `gorm:"column:caller"`
	Type         string    `gorm:"column:type"`
	Name         string    `gorm:"column:name"`
	Payload      string    `gorm:"column:payload"`
	Status       string    `gorm:"column:status"`
	ErrorMessage string    `gorm:"column:error_message"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	FinishedAt   time.Time `gorm:"column:finished_at"`
}

type Repository interface {
	InsertJob(ctx context.Context, req entities.Job) (string, error)
	UpdateJob(ctx context.Context, reqID string, data entities.Job) error
	GetJobStatus(ctx context.Context, reqID string) (string, error)
}

type repository struct {
	DbRepository *mysql.DBRepository
}

func NewNotificationRepository(db *mysql.DBRepository) Repository {
	return &repository{
		DbRepository: db,
	}
}

func (r *repository) UpdateJob(ctx context.Context, reqID string, data entities.Job) error {
	span, ctx := apm.StartSpan(ctx, "repositories/notification/notification.go", "UpdateJob")
	defer span.End()
	d := r.DbRepository.NotificationDB
	err := d.Where("request_id=?", reqID).Updates(data)
	return err.Error
}

func (r *repository) InsertJob(ctx context.Context, req entities.Job) (string, error) {
	span, ctx := apm.StartSpan(ctx, "repositories/notification/notification.go", "InsertJob")
	defer span.End()
	err := req.Insert(r.DbRepository, ctx)
	return req.RequestID, err
}

func (r *repository) GetJobStatus(ctx context.Context, reqID string) (string, error) {
	span, ctx := apm.StartSpan(ctx, "repositories/notification/notification.go", "GetJobStatus")
	defer span.End()
	var status string
	d := r.DbRepository.NotificationDB
	if err := d.Raw("SELECT job.status FROM job WHERE job.request_id = ?", reqID).Scan(&status).Error; err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return status, nil
}
