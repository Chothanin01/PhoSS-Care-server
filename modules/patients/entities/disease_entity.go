package entities

import (

	"github.com/google/uuid"
)

type Disease struct {
	DiseaseID uuid.UUID   `json:"disease_id"`
	Name      string `json:"name"`
}

type GetDiseaseListRes struct {
	Diseases []Disease `json:"diseases"`
}

type DiseaseWithStatus struct {
	DiseaseID uuid.UUID   `json:"disease_id"`
	Name           string `json:"name"`
	HasAppointment bool   `json:"has_appointment"`
}

type DiseaseRepository interface {
	LinkPatientDiseases(patientID uuid.UUID, diseases []PatientDiseaseEntity) error
}

type DiseaseGetRepository interface {
	GetAllDiseases() ([]Disease, error)
}

type DiseaseUsecase interface {
	GetAllDiseases() ([]Disease, error)
}

