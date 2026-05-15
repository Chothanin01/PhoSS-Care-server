package repositories

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
)

type notificationQueryRepo struct {
	db *gorm.DB
}

func NewNotificationQueryRepo(db *gorm.DB) entities.NotificationQueryRepo {
	return &notificationQueryRepo{db: db}
}

type notificationCommandRepo struct {
	db *gorm.DB
}

func NewNotificationCommandRepo(db *gorm.DB) entities.NotificationCommandRepo {
	return &notificationCommandRepo{db: db}
}

func (r *notificationQueryRepo) GetPatientNotifications(patientID uuid.UUID) ([]entities.NotificationItem, error) {
	var results []struct {
		databases.Notification
		DiseaseID *uuid.UUID `gorm:"column:disease_id"` 
	}

	err := r.db.Model(&databases.Notification{}).
		Select("notification.*, appoint.disease_id").
		Joins("LEFT JOIN appoint ON appoint.id = notification.appoint_id").
		Where("notification.patient_id = ?", patientID).
		Order("notification.created_at DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	items := make([]entities.NotificationItem, len(results))
	for i, res := range results {
		items[i] = entities.NotificationItem{
			ID:        res.ID,
			Header:    res.Header,
			Body:      res.Body,
			IsRead:    res.IsRead,
			AppointID: res.AppointID,
			DiseaseID: res.DiseaseID,
			CreatedAt: res.CreatedAt.In(loc).Format("2006-01-02 15:04:05"),
		}
	}

	return items, nil
}

func (r *notificationQueryRepo) CheckHasUnread(patientID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&databases.Notification{}).
		Where("patient_id = ? AND is_read = ?", patientID, false).
		Count(&count).Error

	return count > 0, err
}

func (r *notificationCommandRepo) MarkAsRead(notiID uuid.UUID, patientID uuid.UUID) error {
	
	result := r.db.Model(&databases.Notification{}).
		Where("id = ? AND patient_id = ?", notiID, patientID).
		Update("is_read", true)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}