package entities

import (
	"time"
	"github.com/google/uuid"
)

type RequestEntity struct {
	RequestType string
	Description string
	Date        time.Time
	Status      string
	PatientID   uuid.UUID
	AppointID   *uuid.UUID 
	DiseaseID   *uuid.UUID 
	CreatedBy   uuid.UUID
}

type RequestItem struct {
	Type          string    `json:"type"`           
	DiseaseID     uuid.UUID `json:"disease_id"`     
	DocumentTypes []string  `json:"document_types"` 
}

type CreateDocumentReq struct {
	Requests []RequestItem `json:"requests"`
}

type AvailableRequestOption struct {
	Name      string     `json:"name"`                 
	Type      string     `json:"type"`                 
	Available bool       `json:"available"`            
	DiseaseID *uuid.UUID `json:"disease_id,omitempty"` 
}

type RequestQueryRepo interface {
	CheckHasAnyAppointment(patientID uuid.UUID) (bool, error)
	CheckHasVaccineHistory(patientID uuid.UUID) (bool, error)
	GetLatestCompletedAppoint(patientID uuid.UUID) (*uuid.UUID, *uuid.UUID, error)
}

type RequestQueryUsecase interface {
	GetAvailableDocumentOptions(patientID uuid.UUID) ([]AvailableRequestOption, error)
}

type RequestCommandRepo interface {
	GetLatestAppointID(patientID uuid.UUID, diseaseID uuid.UUID) (*uuid.UUID, error)
	GetDiseaseName(diseaseID uuid.UUID) (string, error)
	
	SaveMultipleRequests(reqs []RequestEntity, notis []NotificationEntity) error
}

type RequestCommandUsecase interface {
	SubmitDocumentRequest(userID uuid.UUID, patientID uuid.UUID, payload *CreateDocumentReq) error
}