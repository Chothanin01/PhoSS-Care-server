package entities

import (
	"time"

	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"github.com/google/uuid"
)

type AppointmentEntity struct {
	ID 		  		uuid.UUID		`json:"appoint_id"`
	No				int				`json:"no"`
	Doctor    		string			`json:"doctor"`
	Status    		string			`json:"status"`
	Purpose   		string			`json:"purpose"`
	Place     		string			`json:"place"`
	Date      		string			`json:"date"`
	StartTime 		string			`json:"start_time"`
	EndTime   		string			`json:"end_time"`
	Symptom    		string			`json:"symptom"`
	Note      		string			`json:"note"`
	Delay	  		bool			`json:"delay"`
	DiseaseID		uuid.UUID 		`json:"disease_id"`
	DiseaseName     string			`json:"disease_name,omitempty"`
	CreatedAt		string			`json:"created_at,omitempty"`
	CreatedBy		string			`json:"created_by,omitempty"`
}

type DelayRequestEntity struct {
	RequestType string
	Description string
	Date        time.Time
	StartTime   string
	EndTime     string
	Status      string
	PatientID   uuid.UUID
	AppointID   uuid.UUID
	DiseaseID   uuid.UUID
	CreatedBy   uuid.UUID
}

type AppointmentDelayReq struct {
	AppointID  uuid.UUID `json:"appoint_id"`
	DiseaseID  uuid.UUID `json:"disease_id"`
	Date       string    `json:"date"` 
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	Description string   `json:"description"`
}

type AppointmentQueryRepo interface {
    GetAppointmentDetail(patientID uuid.UUID, diseaseID uuid.UUID) (*databases.Appoint, *databases.Admin, error)
    
    ListPatientAppointments(patientID uuid.UUID) ([]databases.Appoint, error)
}

type AppointmentQueryUsecase interface {
    GetAppointmentDetail(patientID uuid.UUID, diseaseID uuid.UUID) (*AppointmentEntity, error)

	ListPatientAppointments(patientID uuid.UUID) (*PatientAppointment, error)
}

type AppointmentCommandRepo interface {
	CheckAppointmentExists(appointID uuid.UUID, patientID uuid.UUID) (bool, error)
	SaveDelayRequest(req *DelayRequestEntity) error
}

type AppointmentCommandUsecase interface {
    SubmitDelayRequest(userID uuid.UUID, patientID uuid.UUID, payload *AppointmentDelayReq) error
    
}