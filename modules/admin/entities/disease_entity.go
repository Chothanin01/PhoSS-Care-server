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
	Disease      []Disease              `json:"diseases"`
}

type AppointmentInfo struct {
	No 	 	int            `json:"no"`
	Date    time.Time      `json:"date"`
	Time    string         `json:"time"`
	Symptom string         `json:"symptom"`
	Note    string         `json:"note"`
	Place   string         `json:"place"`
	Doctor  string         `json:"doctor"`
	Status  string         `json:"status"`
	Letter  bool           `json:"letter"`
	Delay   bool           `json:"delay"`
	Health  *HealthInfo    `json:"health,omitempty"`
}

type HealthInfo struct {
	Height float32 `json:"height"`
	Weight float32 `json:"weight"`
	BMI    float64 `json:"bmi"`
	Pulse  int     `json:"pulse"`
	Sugar  int     `json:"sugar"`
}

type DiseaseRepo interface {
	LinkPatientDiseases(patientID uuid.UUID, diseases []PatientDiseaseEntity) error
}

type DiseaseGetRepo interface {
	GetAllDiseases() ([]Disease, error)
}

type DiseaseUsecase interface {
	GetAllDiseases() ([]Disease, error)
}

