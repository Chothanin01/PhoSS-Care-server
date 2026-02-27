package entities

import (

	"github.com/google/uuid"
)

type RequestInfo struct {
	ID           uuid.UUID `json:"id"`
	RequestType  string    `json:"req_type"`
	PatientName  string    `json:"patient_name"`
	DiseaseName  string    `json:"disease_name"`
	HnNumber     string    `json:"hn_number"`
	Status       string    `json:"status"`
	Description  string    `json:"description"`
	Date		 string    `json:"date"`
	Time		 string    `json:"time"`
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
}

type RequestGetRepo interface {
	GetRequestInfoByID(id uuid.UUID) (*RequestInfoRes, error)
	GetRequestsWithFilter(params RequestQueryParams) ([]RequestInfo, error)
	CountRequestsWithFilter(params RequestQueryParams) (int64, error)
}

type RequestGetUsecase interface {
	GetRequestsWithFilter(params RequestQueryParams) (*RequestListRes, error)
	GetRequestInfoByID(id uuid.UUID) (*RequestInfoRes, error)
}