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

type DiseaseScheduleEntity struct {
	DiseaseID     uuid.UUID `json:"disease_id"`
	DiseaseName   string    `json:"disease_name"`
	AvailableDays []string  `json:"available_days"` 
}

type HistoryAppointEntity struct {
	AppointID   uuid.UUID `json:"appoint_id"`
	No          int       `json:"no"`
	Date        string    `json:"date"`
	Note        string    `json:"note"`
	ColorStatus string    `json:"color_status"` 
	DoctorName  string    `json:"doctor"`
}

type DiseaseHistoryResponse struct {
	TotalPages   int                 `json:"total_pages"`
	CurrentPage  int                 `json:"current_page"`
	Appointments []HistoryAppointEntity `json:"appointments"`
}

type HealthEntity struct {
	Pulse    int     `json:"pulse"`
	Pressure int     `json:"pressure"`
	Height   int     `json:"height"`
	Weight   float64 `json:"weight"`
	BMI      float64 `json:"bmi"`
}

type HistoryDetailEntity struct {
	AppointID   uuid.UUID    `json:"appoint_id"`
	No          int          `json:"no"`
	Date        string       `json:"date"`
	Note        string       `json:"note"`
	ColorStatus string       `json:"color_status"`
	Doctor      string       `json:"doctor"`
	Purpose     string       `json:"purpose"`
	Symptom     string       `json:"symptom"`
	Health      HealthEntity `json:"health"`
	NextAppointID *uuid.UUID `json:"next_appoint_id"`
	PrevAppointID *uuid.UUID `json:"prev_appoint_id"`
}

type AppointmentQueryRepo interface {
	CheckPatientHasDisease(patientID uuid.UUID, diseaseID uuid.UUID) (bool, error)
    GetAppointmentDetail(patientID uuid.UUID, diseaseID uuid.UUID) (*databases.Appoint, *databases.Admin, error)
	GetScheduleByDisease(diseaseID uuid.UUID) (*DiseaseScheduleEntity, error)
    ListPatientAppointments(patientID uuid.UUID) ([]databases.Appoint, error)

	CountDiseaseHistory(patientID uuid.UUID, diseaseID uuid.UUID) (int64, error)
	GetDiseaseHistory(patientID uuid.UUID, diseaseID uuid.UUID, limit int, offset int) ([]HistoryAppointEntity, error)
	GetHistoryDetail(appointID uuid.UUID, patientID uuid.UUID) (*HistoryDetailEntity, error)
}


type AppointmentQueryUsecase interface {
    GetAppointmentDetail(patientID uuid.UUID, diseaseID uuid.UUID) (*AppointmentEntity, error)
	GetScheduleByDisease(patientID uuid.UUID, diseaseID uuid.UUID) (*DiseaseScheduleEntity, error)
	ListPatientAppointments(patientID uuid.UUID) (*PatientAppointment, error)

	GetDiseaseHistory(patientID uuid.UUID, diseaseID uuid.UUID, page int) (*DiseaseHistoryResponse, error)
	GetHistoryDetail(appointID uuid.UUID, patientID uuid.UUID) (*HistoryDetailEntity, error)
}

type AppointmentCommandRepo interface {
	CheckAppointmentExists(appointID uuid.UUID, patientID uuid.UUID) (bool, error)
	SaveDelayRequest(req *DelayRequestEntity) error
}

type AppointmentCommandUsecase interface {
    SubmitDelayRequest(userID uuid.UUID, patientID uuid.UUID, payload *AppointmentDelayReq) error
    
}