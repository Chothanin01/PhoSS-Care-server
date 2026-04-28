package repositories

import (
	
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
	var dbNotis []databases.Notification

	err := r.db.Where("patient_id = ?", patientID).
		Order("created_at DESC").
		Find(&dbNotis).Error

	if err != nil {
		return nil, err
	}

	items := make([]entities.NotificationItem, len(dbNotis))
	for i, noti := range dbNotis {
		items[i] = entities.NotificationItem{
			ID:        noti.ID,
			Header:    noti.Header,
			Body:      noti.Body,
			IsRead:    noti.IsRead,
			AppointID: noti.AppointID, 
			CreatedAt: noti.CreatedAt.Format("2006-01-02 15:04:05"),
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