package entities

import (
	"time"

	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"github.com/google/uuid"
)

// --- Requests ---

type AppointmentCreateReq struct {
	PatientID uuid.UUID `json:"patient_id"`
	DiseaseID uuid.UUID `json:"disease_id"`

	DoctorID     uuid.UUID `json:"doctor_id"`
	NextDoctorID uuid.UUID `json:"next_doctor_id"`

	Purpose   string 	`json:"purpose"`
	Prepare   string   	`json:"prepare"`
	Place     string   	`json:"place"`
	Date      string   	`json:"date"`
	StartTime string   	`json:"start_time"`
	EndTime   string   	`json:"end_time"`
	Symptom   string   	`json:"symptom"`
	Note      string   	`json:"note"`
	Health    Health   	`json:"health"`
}

type AppointmentUpdateReq struct {
	AppointID uuid.UUID `json:"appoint_id"`
	DoctorID  uuid.UUID `json:"doctor_id"`

	Purpose   string `json:"purpose"`
	Place     string `json:"place"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type VaccineAppointmentCreateReq struct {
	PatientID    uuid.UUID `json:"patient_id"`
	OldVaccineID uuid.UUID `json:"old_vaccine_id"`
	VaccineID    uuid.UUID `json:"vaccine_id"`
	DoseNumber   int       `json:"dose_number"`

	VaccineDoctorID uuid.UUID `json:"vaccine_doctor_id"`
	DoctorID        uuid.UUID `json:"doctor_id"`

	Place     string `json:"place"`
	Date      string `json:"date"`      
	NextDate  string `json:"next_date"` 
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type VaccineAppointmentUpdateReq struct {
	AppointID uuid.UUID `json:"appoint_id"`
	VaccineID uuid.UUID `json:"vaccine_id"`
	PatientID uuid.UUID `json:"patient_id"`
	DoctorID  uuid.UUID `json:"doctor_id"`

	Place     string `json:"place"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// --- Domain Entities ---

type AppointmentEntity struct {
	ID          uuid.UUID
	Status      string
	Purpose     string
	Place       string
	Date        string
	StartTime   string
	EndTime     string
	Symptom     string
	Note        string
	ColorStatus string
	Health      Health
	PatientID   uuid.UUID
	DiseaseID   uuid.UUID
	DoctorID    uuid.UUID
	CreatedBy   uuid.UUID
	UpdatedBy   uuid.UUID
}

type Health struct {
	Weight   float64 `json:"weight"`
	Height   int     `json:"height"`
	Pulse    int     `json:"pulse"`
	Sugar    int     `json:"sugar"`
	BMI      float64 `json:"bmi"`
	Pressure int     `json:"pressure"`
}

type AppointmentRes struct {
	ID          uuid.UUID `json:"id"`
	Status      string    `json:"status"`
	No          int       `json:"no"`
	Date        time.Time `json:"date"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	DoctorID    uuid.UUID `json:"doctor_id"` 
	Purpose     string    `json:"purpose"`
	Place       string    `json:"place"`
	ColorStatus string    `json:"color_status"`
}

type VaccineAppointmentRes struct {
	AppointID uuid.UUID `json:"appoint_id"`
	No        int       `json:"no"`
	Status    string    `json:"status"`
	Date      string    `json:"date"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	DoctorID  uuid.UUID `json:"doctor_id"`
}


type AppointmentUsecase interface {
	CreateAppointment(req *AppointmentCreateReq, adminID uuid.UUID) (*AppointmentRes, error)
	FindOngoingVaccination(patientID uuid.UUID) (*VaccineFullDetail, error)
	CreateVaccineAppointment(req *VaccineAppointmentCreateReq, adminID uuid.UUID) (*VaccineAppointmentRes, error)
	UpdateAppointment(req *AppointmentUpdateReq, adminID uuid.UUID) (*AppointmentRes, error)
	UpdateVaccineAppointment(req *VaccineAppointmentUpdateReq, adminID uuid.UUID) (*VaccineAppointmentRes, error)
}

type AppointmentRepository interface {
	FindOngoing(patientID, diseaseID uuid.UUID) (*databases.Appoint, error)
	FindByID(appointID uuid.UUID) (*databases.Appoint, error)
	FindOngoingVaccination(patientID uuid.UUID) (*databases.VaccinationRecord, error)
	FindPatientVaccine(patientID uuid.UUID) (bool, error)
	CheckVaccineExists(vaccineID uuid.UUID) (bool, error)
	FindVaccineDiseaseID() (uuid.UUID, error)
	IsVaccineDisease(diseaseID uuid.UUID) (bool, error)
	FindMaxDose(patientID, vaccineID uuid.UUID) (int, error)

	UpdateSymptomNote(doctorID uuid.UUID, appointID uuid.UUID, symptom string, note string, adminID uuid.UUID) error
	UpdateVaccinationRecord(appointID uuid.UUID, vaccineID uuid.UUID, dose int, adminID uuid.UUID) error
	CompleteAppoint(appointID uuid.UUID, adminID uuid.UUID) error
	CreateAppointment(entity *AppointmentEntity, adminID uuid.UUID) (*databases.Appoint, error)
	DiseaseExists(diseaseID uuid.UUID) (bool, error)

	CreateVaccinationRecord(vaccineID uuid.UUID, appointID uuid.UUID, dose int, doctorID uuid.UUID, adminID uuid.UUID) error
	UpdateVaccineDoctor(recordID uuid.UUID, doctorID uuid.UUID, adminID uuid.UUID) error
	CompleteVaccinationRecord(recordID uuid.UUID, adminID uuid.UUID) error

	UpdateHealth(appointID uuid.UUID, health *Health, adminID uuid.UUID) error
	UpdateAppointment(appointID uuid.UUID, purpose string, place string, date string, startTime string, endTime string, doctorID uuid.UUID, adminID uuid.UUID) (*databases.Appoint, error)
	UpdateVaccineAppointment(appointID uuid.UUID, place string, date string, start string, end string, doctorID uuid.UUID, adminID uuid.UUID) (*databases.Appoint, error)
	UpdatePatientHealth(patientID uuid.UUID, weight float64, height int, adminID uuid.UUID) error
}

type AppointmentTransaction interface {

    Do(fn func(RepositorySet) error) error

}