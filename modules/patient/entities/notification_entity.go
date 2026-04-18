package entities

import (

	"github.com/google/uuid"
)

const (
	NotiHeaderMedical  = "ใบรับรองแพทย์"
	NotiHeaderDocument = "เอกสารรับรอง"
	NotiHeaderDelay    = "การเลื่อนนัด"

	NotiBodyMedicalSent  = "ระบบได้ยื่นคำขอใบรับรองแพทย์ไปแล้ว"
	NotiBodyDocumentSent = "ระบบได้ยื่นคำขอเอกสารรับรองไปแล้ว"
	
	NotiBodyMedicalSuccess = "โรงพยาบาลได้เตรียมเอกสารของคุณเเล้ว"
	NotiBodyDelaySuccess   = "ระบบได้ยืนยันการเลื่อนนัดของคุณแล้ว"
	NotiBodyDelayDenied    = "ระบบได้ปฏิเสธการเลื่อนนัดของคุณ กรุณาเลื่อนนัดใหม่อีกครั้ง"
)

type NotificationEntity struct {
	Header    string
	Body      string
	PatientID uuid.UUID
	CreatedBy uuid.UUID
}

type NotificationItem struct {
	ID        uuid.UUID  `json:"id"`
	Header    string     `json:"header"`
	Body      string     `json:"body"`
	IsRead    bool       `json:"is_read"`
	AppointID *uuid.UUID `json:"appoint_id,omitempty"`
	CreatedAt string     `json:"created_at"`
}

type NotificationQueryRepo interface {
	GetPatientNotifications(patientID uuid.UUID) ([]NotificationItem, error)
	CheckHasUnread(patientID uuid.UUID) (bool, error)
}

type NotificationQueryUsecase interface {
	GetPatientNotifications(patientID uuid.UUID) ([]NotificationItem, error)
	CheckHasUnread(patientID uuid.UUID) (bool, error)
}

type NotificationCommandRepo interface {
	MarkAsRead(notiID uuid.UUID, patientID uuid.UUID) error
}

type NotificationCommandUsecase interface {
	MarkAsRead(notiID uuid.UUID, patientID uuid.UUID) error
}