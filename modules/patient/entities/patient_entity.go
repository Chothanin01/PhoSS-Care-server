package entities

import (

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases" 
)

type PatientBasicInfoRes struct {
	PatientID uuid.UUID `json:"patient_id"`
	FullName  string    `json:"fullname"`
	HnNumber  string    `json:"hn_number"`
	AgeYears  int       `json:"age_years"`
	AgeMonths int       `json:"age_months"`
	AgeDays   int       `json:"age_days"`
}

type PatientGetRepo interface {
	GetPatientBasicInfo(userID uuid.UUID) (*databases.Patient, error)
}

type PatientGetUsecase interface {
	GetPatientBasicInfo(userID uuid.UUID) (*PatientBasicInfoRes, error)
}