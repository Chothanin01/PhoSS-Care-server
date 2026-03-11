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
    StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`

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
	StartTime string
	EndTime   string
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

type AppointmentRes struct {
	ID      uuid.UUID `json:"id"`
	Status  string    `json:"status"`
	No      int       `json:"no"`
	Date    time.Time `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Doctor  string    `json:"doctor"`
	Purpose string    `json:"purpose"`
	Place   string    `json:"place"`
}

type AppointmentUpdateReq struct {
	AppointID uuid.UUID `json:"appoint_id"`

	Purpose string `json:"purpose"`
	Place   string `json:"place"`

	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`

	DoctorTitle     string `json:"doctor_title"`
	DoctorFirstName string `json:"doctor_firstname"`
	DoctorLastName  string `json:"doctor_lastname"`
}

type VaccineAppointmentCreateReq struct {
	PatientID uuid.UUID `json:"patient_id"`

	OldVaccineID uuid.UUID `json:"old_vaccine_id"`
	VaccineID    uuid.UUID `json:"vaccine_id"`

	DoseNumber int `json:"dose_number"`
	NextDoseNumber int `json:"next_dose_number"`

	VaccineDoctorTitle     string `json:"vaccine_doctor_title"`
	VaccineDoctorFirstName string `json:"vaccine_doctor_firstname"`
	VaccineDoctorLastName  string `json:"vaccine_doctor_lastname"`

	DoctorTitle     string `json:"doctor_title"`
	DoctorFirstName string `json:"doctor_firstname"`
	DoctorLastName  string `json:"doctor_lastname"`
	Place  string `json:"place"`

	Date string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type VaccineAppointmentRes struct {
	AppointID uuid.UUID `json:"appoint_id"`
	No        int       `json:"no"`
	Status    string    `json:"status"`
	Date      string    `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Doctor    string    `json:"doctor"`
}

type VaccineAppointmentUpdateReq struct {
	AppointID uuid.UUID `json:"appoint_id"`
	VaccineID uuid.UUID `json:"vaccine_id"`
	PatientID uuid.UUID `json:"patient_id"`

	Place string `json:"place"`
	Date  string `json:"date"`

	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`

	DoctorTitle     string `json:"doctor_title"`
	DoctorFirstName string `json:"doctor_firstname"`
	DoctorLastName  string `json:"doctor_lastname"`
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
	

	UpdateSymptomNote(doctor string, appointID uuid.UUID, symptom string, note string, adminID uuid.UUID) error
	UpdateVaccinationRecord(appointID uuid.UUID, vaccineID uuid.UUID, dose int, adminID uuid.UUID) (error)
	CompleteAppoint(appointID uuid.UUID, adminID uuid.UUID) error
	CreateAppointment(entity *AppointmentEntity, adminID uuid.UUID) (*databases.Appoint, error)
	CreateHealthRecord(health *Health, patientID uuid.UUID, appointID uuid.UUID, adminID uuid.UUID) error

	DiseaseExists(diseaseID uuid.UUID) (bool, error)

	CreateVaccinationRecord(vaccineID uuid.UUID, appointID uuid.UUID, dose int, doctor string, adminID uuid.UUID,) error
	UpdateVaccineDoctor(recordID uuid.UUID, doctor string, adminID uuid.UUID) error
	CompleteVaccinationRecord(recordID uuid.UUID, adminID uuid.UUID) error

	UpdateAppointment( appointID uuid.UUID, purpose string, place string, date string, startTime string, endTime string, doctor string, adminID uuid.UUID) (*databases.Appoint, error)
	UpdateVaccineAppointment(appointID uuid.UUID, place string, date string, start string, end string, doctor string, adminID uuid.UUID) (*databases.Appoint, error)
}

type AppointmentTransaction interface {
	Do(fn func(RepositorySet) error) error
}
