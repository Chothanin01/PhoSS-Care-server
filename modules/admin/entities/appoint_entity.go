package entities

import (
	"time"

	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"github.com/google/uuid"
)

type AppointmentCreateReq struct {
    PatientID uuid.UUID `json:"patient_id"`
    DiseaseID uuid.UUID `json:"disease_id"`

    DoctorTitle     string `json:"doctor_title"`
    DoctorFirstName string `json:"doctor_firstname"`
    DoctorLastName  string `json:"doctor_lastname"`

    NextDoctorTitle     string `json:"next_doctor_title"`
    NextDoctorFirstName string `json:"next_doctor_firstname"`
    NextDoctorLastName  string `json:"next_doctor_lastname"`

    Purpose string `json:"purpose"`
    Place   string `json:"place"`
    Date    string `json:"date"`
    Time    string `json:"time"`

    Symptom string `json:"symptom"`
    Note    string `json:"note"`

    Health Health `json:"health"`
}

type AppointmentEntity struct {
	Doctor    string
	Status    string
	Purpose   string
	Place     string
	Date      string
	Time      string
	Symtom    string
	Note      string
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
	ID      uuid.UUID `json:"id"`
	Status  string    `json:"status"`
	No      int       `json:"no"`
	Date    time.Time `json:"date"`
	Time    string    `json:"time"`
	Doctor  string    `json:"doctor"`
	Purpose string    `json:"purpose"`
	Place   string    `json:"place"`
}

type AppointmentUpdateReq struct {
	AppointID       uuid.UUID `json:"appoint_id"`
	DoctorTitle     string    `json:"doctor_title"`
	DoctorFirstName string    `json:"doctor_firstname"`
	DoctorLastName  string    `json:"doctor_lastname"`
	Symptom         string    `json:"symptom"`
	Purpose         string    `json:"purpose"`
	Place           string    `json:"place"`
	Date            string    `json:"date"`
	Time            string    `json:"time"`
	Health          Health    `json:"health"`
}

type AppointmentUpdateRes struct {
	ID      uuid.UUID `json:"id"`
	Doctor  string    `json:"doctor"`
	Status  string    `json:"status"`
	Purpose string    `json:"purpose"`
	Place   string    `json:"place"`
	Date    string    `json:"date"`
	Time    string    `json:"time"`
}

type AppointmentUsecase interface {
	CreateAppointment(req *AppointmentCreateReq, adminID uuid.UUID) (*AppointmentCreateRes, error)
	FindOngoingVaccination(patientID uuid.UUID) (*VaccineFullDetail, error) 
}

type AppointmentRepository interface {
	FindOngoing(patientID, diseaseID uuid.UUID) (*databases.Appoint, error)
	FindByID(appointID uuid.UUID) (*databases.Appoint, error)
	FindOngoingVaccination(patientID uuid.UUID) (*databases.VaccinationRecord, error)

	UpdateSymptomNote(doctor string, appointID uuid.UUID, symptom string, note string, adminID uuid.UUID) error
	CompleteAppoint(appointID uuid.UUID, adminID uuid.UUID) error
	CreateAppointment(entity *AppointmentEntity, adminID uuid.UUID) (*databases.Appoint, error)
	CreateHealthRecord(health *Health, patientID uuid.UUID, appointID uuid.UUID, adminID uuid.UUID) error

	DiseaseExists(diseaseID uuid.UUID) (bool, error)
}

type AppointmentTransaction interface {
	Do(fn func(RepositorySet) error) error
}
