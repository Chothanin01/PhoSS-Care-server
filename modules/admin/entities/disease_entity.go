package entities

import (
	"time"

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
	Health  	*HealthInfo    `json:"health,omitempty"`
}

type HealthInfo struct {
	Height float32 `json:"height"`
	Weight float32 `json:"weight"`
	BMI    float64 `json:"bmi"`
	Pulse  int     `json:"pulse"`
	Sugar  int     `json:"sugar"`
}

type VaccineFullDetail struct {
	VaccineID uuid.UUID `json:"vaccine_id"`
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

