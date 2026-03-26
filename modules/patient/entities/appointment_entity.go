package entities

import (

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

type GetAppointmentRepo interface {
	GetAppointmentDetail(patientID uuid.UUID, diseaseID uuid.UUID) (*databases.Appoint, *databases.Admin, error)
}

type GetAppointmentUsecase interface {
	GetAppointmentDetail(patientID uuid.UUID, diseaseID uuid.UUID) (*AppointmentEntity, error)
}