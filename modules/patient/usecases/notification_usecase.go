package usecases
import (
	"fmt"

 	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
)

type notificationQueryUsecase struct { 
	repo entities.NotificationQueryRepo 
}
func NewNotificationQueryUsecase(repo entities.NotificationQueryRepo) entities.NotificationQueryUsecase { 
	return &notificationQueryUsecase{repo: repo} 
}

type notificationCommandUsecase struct { 
	repo entities.NotificationCommandRepo 
}

func NewNotificationCommandUsecase(repo entities.NotificationCommandRepo) entities.NotificationCommandUsecase { 
	return &notificationCommandUsecase{repo: repo} 
}

func (u *notificationQueryUsecase) GetPatientNotifications(patientID uuid.UUID) ([]entities.NotificationItem, error) {
	
	if patientID == uuid.Nil {
		return nil, fmt.Errorf("patient ID is required to fetch notifications")
	}

	notifications, err := u.repo.GetPatientNotifications(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}

	return notifications, nil
}

func (u *notificationQueryUsecase) CheckHasUnread(patientID uuid.UUID) (bool, error) {
	
	if patientID == uuid.Nil {
		return false, fmt.Errorf("patient ID is required to check unread status")
	}

	hasUnread, err := u.repo.CheckHasUnread(patientID)
	if err != nil {
		return false, fmt.Errorf("failed to check unread status: %w", err)
	}

	return hasUnread, nil
}

func (u *notificationCommandUsecase) MarkAsRead(notiID uuid.UUID, patientID uuid.UUID) error {
	
	if notiID == uuid.Nil {
		return fmt.Errorf("notification ID is required")
	}
	if patientID == uuid.Nil {
		return fmt.Errorf("patient ID is required to verify ownership")
	}

	err := u.repo.MarkAsRead(notiID, patientID)
	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}