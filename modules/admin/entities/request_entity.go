package entities

import (
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"github.com/google/uuid"
)

type RequestInfo struct {
	ID           uuid.UUID  `json:"id"`
	RequestType  string     `json:"req_type"`
	PatientName  string     `json:"patient_name"`
	DiseaseName  string     `json:"disease_name"`
	HnNumber     string     `json:"hn_number"`
	Status       string     `json:"status"`
	Description  string     `json:"description"`
	Date		 string     `json:"date"`
	Time		 string     `json:"time"`
	AppointID	 uuid.UUID  `json:"appoint_id"`
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
	RequestID        uuid.UUID  `json:"request_id"`
	RequestType      string     `json:"request_type"`
	Status           string     `json:"status"`
	PatientID        uuid.UUID  `json:"patient_id"`
	FullName         string     `json:"fullname"`
	Date             string     `json:"date,omitempty"`
	Time             string     `json:"time,omitempty"`
	CreatedDate      string     `json:"created_date,omitempty"`
	Description      string     `json:"description,omitempty"`
	IDCard           string     `json:"id_card"`
	HnNumber         string     `json:"hn_number"`
	Doctor           string     `json:"doctor,omitempty"`
	AppointDate      string     `json:"appoint_date,omitempty"`
	AppointTime      string     `json:"appoint_time,omitempty"`
	DiseaseName  	 string     `json:"disease_name"`
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
	UpdateRequestStatus(id uuid.UUID, status, description string, adminID uuid.UUID) error
	UpdateAppointForAccepted(appointID uuid.UUID, adminID uuid.UUID) error
}

type RequestUpdateUsecase interface {
	UpdateRequestStatus(req *RequestStatusUpdateReq, adminID uuid.UUID) (*RequestStatusUpdateRes, error)
}
