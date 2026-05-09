package entities

import (
	"time"

	"github.com/google/uuid"
)

type Disease struct {
	DiseaseID uuid.UUID   `json:"disease_id"`
	Name      string `json:"name"`
	AppointID *uuid.UUID `json:"appoint_id,omitempty"`
	AppointDate *time.Time `json:"appoint_date,omitempty"`
}

type GetDiseaseListRes struct {
	Diseases []Disease `json:"diseases"`
}

type DiseaseWithStatus struct {
	DiseaseID uuid.UUID   `json:"disease_id"`
	Name           string `json:"name"`
	HasAppointment bool   `json:"has_appointment"`
	HasOverdue     bool      `json:"has_overdue"`
}

type DiseaseInfoRes struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []DiseaseData `json:"data"`
}

type DiseaseData struct {
	PatientID    uuid.UUID      		`json:"patient_id"`
	Fullname     string         		`json:"fullname"`
	DiseaseID    uuid.UUID              `json:"disease_id"`
	DiseaseName  string               	`json:"disease_name"`
	Appointment  []AppointmentInfo 		`json:"appointment_info"`
}

type AppointmentInfo struct {
	No 	 		int            `json:"no"`
	Date    	time.Time      `json:"date"`
	StartTime 	string 		    `json:"start_time"`
	EndTime   	string 		   `json:"end_time"`
	Symptom 	string         `json:"symptom"`
	Note    	string         `json:"note"`
	Place   	string         `json:"place"`
	Doctor  	string         `json:"doctor"`
	Status  	string         `json:"status"`
	Letter  	bool           `json:"letter"`
	Delay   	bool           `json:"delay"`
	ColorStatus string 		   `json:"color_status"`
	Health  	*Health    `json:"health,omitempty"`
}

type VaccineFullDetail struct {
	VaccineID uuid.UUID `json:"vaccine_id"`
	Name 	  string    `json:"name"`
	Date      string    `json:"date"`
	Type      string    `json:"type"`
	Effect    string    `json:"effect"`
	Note      string    `json:"note"`
	Age       string    `json:"age"`
}

type VaccinationRecordEntity struct {
	ID        uuid.UUID
	PatientID uuid.UUID
	VaccineID uuid.UUID
	AppointID uuid.UUID
	Dose      int
	Status    string
}

type DiseaseRepo interface {
	LinkPatientDiseases(patientID uuid.UUID, diseases []PatientDiseaseEntity) error
}

type DiseaseGetRepo interface {
	GetAllDiseases() ([]Disease, error)
	GetAllVaccines() ([]VaccineFullDetail, error)
}

type DiseaseUsecase interface {
	GetAllDiseases() ([]Disease, error)
	GetAllVaccines() ([]VaccineFullDetail, error)
}

