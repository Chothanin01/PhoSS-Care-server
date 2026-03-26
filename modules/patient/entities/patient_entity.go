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

type PatientAppointment struct {
	PatientID 		uuid.UUID 				`json:"patient_id"`
	FullName  		string    				`json:"fullname"`
	HnNumber  		string    				`json:"hn_number"`
	AgeYears  		int       				`json:"age_years"`
	AgeMonths 		int       				`json:"age_months"`
	AgeDays   		int       				`json:"age_days"`
	Appointments 	[]AppointmentEntity 	`json:"appoint"`
}


type GetPatientRepo interface {
	GetPatientBasicInfo(patientID uuid.UUID) (*databases.Patient, error)
	GetPatientAppointment(patientID uuid.UUID) ([]databases.Appoint, error)
}

type GetPatientUsecase interface {
	GetPatientBasicInfo(patientID uuid.UUID) (*PatientBasicInfoRes, error)
	GetPatientAppointment(patientID uuid.UUID) (*PatientAppointment, error)
}