package entities

import (
	"time"
	
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"github.com/google/uuid"
)

type RequestInfo struct {
	ID           uuid.UUID  `json:"id"`
	RequestType  string     `json:"req_type"`
	PatientName  string     `json:"patient_name"`
	DiseaseID    uuid.UUID `json:"disease_id"`
	DiseaseName  string    `json:"disease_name"`
	HnNumber     string     `json:"hn_number"`
	Status       string     `json:"status"`
	Description  string     `json:"description"`
	Date		 time.Time  `json:"date"`
	StartTime 	 string  	`json:"start_time"`
	EndTime   	 string 	`json:"end_time"`
	AppointID	 uuid.UUID  `json:"appoint_id"`
	PatientID    uuid.UUID  `json:"patient_id"`
}

type RequestQueryParams struct {
	ReqType string `query:"req_type"`
	Status  string `query:"status"`
	Page    int    `query:"page"`
	Limit   int    `query:"limit"`
}

type RequestListRes struct {
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
	Data       []RequestInfo  `json:"data"`
}

type RequestInfoRes struct {
	RequestID        		uuid.UUID  `json:"request_id"`
	RequestType      		string     `json:"request_type"`
	Status           		string     `json:"status"`
	PatientID        		uuid.UUID  `json:"patient_id"`
	FullName         		string     `json:"fullname"`
	Date             		string     `json:"date,omitempty"`
	StartTime 		 		string     `json:"start_time"`
	EndTime   		 		string 	   `json:"end_time"`
	CreatedDate      		string     `json:"created_date,omitempty"`
	Description      		string     `json:"description,omitempty"`
	IDCard           		string     `json:"id_card"`
	HnNumber         		string     `json:"hn_number"`
	Doctor           		string     `json:"doctor,omitempty"`
	AppointID      	  		uuid.UUID  `json:"appoint_id,omitempty"`
	AppointDate      	  	string     `json:"appoint_date,omitempty"`
	AppointStartTime      	string     `json:"appoint_start_time,omitempty"`
	AppointEndTime   		string     `json:"appoint_end_time,omitempty"`
	DiseaseName  	 		string     `json:"disease_name"`
}

type RequestStatusUpdateReq struct {
	RequestID   uuid.UUID `json:"request_id"`
	Status      string    `json:"status"`       
	Description string    `json:"description"`  
}

type RequestStatusUpdateRes struct {
	ID          uuid.UUID `json:"id"`
	RequestType string    `json:"request_type"`
	Status      string    `json:"status"`
}

type NotificationEntity struct {
	Header    string
	Body      string
	PatientID uuid.UUID
	CreatedBy uuid.UUID
	RequestID *uuid.UUID 
	AppointID *uuid.UUID
}

type RequestGetRepo interface {
	GetRequestInfoByID(id uuid.UUID) (*databases.Request, error)
	GetRequestsWithFilter(params RequestQueryParams) ([]RequestInfo, error)
	CountRequestsWithFilter(params RequestQueryParams) (int64, error)
}

type RequestGetUsecase interface {
	GetRequestsWithFilter(params RequestQueryParams) (*RequestListRes, error)
	GetRequestInfoByID(id uuid.UUID) (*RequestInfoRes, error)
}

type RequestUpdateRepo interface {
	FindRequestByID(id uuid.UUID) (*RequestInfo, error)
	UpdateRequestStatus(id uuid.UUID, status, description string, adminID uuid.UUID, noti NotificationEntity) error
	UpdateAppointForAccepted(requestID uuid.UUID, appointID uuid.UUID, date time.Time, startTime string, endTime string, adminID uuid.UUID, noti NotificationEntity) error 
}

type RequestUpdateUsecase interface {
	UpdateRequestStatus(req *RequestStatusUpdateReq, adminID uuid.UUID) (*RequestStatusUpdateRes, error)
}
