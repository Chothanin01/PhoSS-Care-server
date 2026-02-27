package entities

import (
	"time"
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
)

type AppointmentCreateReq struct {
	PatientID        uuid.UUID `json:"patient_id"`
	DiseaseID        uuid.UUID `json:"disease_id"`
	DoctorTitle      string    `json:"doctor_title"`
	DoctorFirstName  string    `json:"doctor_firstname"`
	DoctorLastName   string    `json:"doctor_lastname"`
	Symptom          string    `json:"symptom"`
	Note             string    `json:"note"`
	Place            string    `json:"place"`
    Time             string    `json:"time"`
    Date             string    `json:"date"`
	Health           Health    `json:"health"`
}

type AppointmentEntity struct {
	Doctor    string
	Status    string
	Note      string
	Place     string
    Date      string
    Time      string
	PatientID uuid.UUID
	DiseaseID uuid.UUID
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
}

type Health struct {
	Weight float64 `json:"weight"`
	Height int     `json:"height"`
	Pulse  int     `json:"pulse"`
	Sugar  int     `json:"sugar"`
	BMI    float64 `json:"bmi"`
}

type AppointmentCreateRes struct {
	ID        uuid.UUID `json:"id"`
	Status    string    `json:"status"`
    No        int       `json:"no"`
	Date      time.Time `json:"date"`
	Time      string    `json:"time"`
	Doctor    string    `json:"doctor"`
	Note      string    `json:"note"`
	Place     string    `json:"place"`
}

type AppointmentUpdateReq struct {
	AppointID       uuid.UUID `json:"appoint_id"`
	DoctorTitle     string    `json:"doctor_title"`
	DoctorFirstName string    `json:"doctor_firstname"`
	DoctorLastName  string    `json:"doctor_lastname"`
	Symptom         string    `json:"symptom"`
	Note            string    `json:"note"`
	Place           string    `json:"place"`
	Date            string    `json:"date"`
	Time            string    `json:"time"`
	Health          Health    `json:"health"`
}
 
type AppointmentUpdateRes struct {
    ID        uuid.UUID `json:"id"`
	Doctor    string    `json:"doctor"`
	Status    string    `json:"status"`
	Note      string    `json:"note"`
	Place     string    `json:"place"`
	Date      string    `json:"date"`
	Time      string    `json:"time"`
}

type AppointmentUsecase interface {
	CreateAppointment(req *AppointmentCreateReq, adminID uuid.UUID) (*AppointmentCreateRes, error)
    UpdateAppointment(req *AppointmentUpdateReq, adminID uuid.UUID) (*AppointmentUpdateRes, error)
}

type AppointmentRepository interface {
	FindOngoing(patientID, diseaseID uuid.UUID) (*databases.Appoint, error)
    FindByID(appointID uuid.UUID) (*databases.Appoint, error)
    UpdateAppointment(e *databases.Appoint) error
    UpdateHealth(appointID uuid.UUID, health *Health, adminID uuid.UUID) error
	CompleteAppoint(appointID uuid.UUID, adminID uuid.UUID) error
	CreateAppointment(entity *AppointmentEntity, adminID uuid.UUID) (*databases.Appoint, error)
	CreateHealthRecord(health *Health, patientID uuid.UUID, appointID uuid.UUID, adminID uuid.UUID) error
    DiseaseExists(diseaseID uuid.UUID) (bool, error)
}

type AppointmentTransaction interface {
	Do(fn func(RepositorySet) error) error
}