package entities

import "github.com/google/uuid"

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

type RequestGetRepo interface {
	GetRequestsWithFilter(params RequestQueryParams) ([]RequestInfo, error)
	CountRequestsWithFilter(params RequestQueryParams) (int64, error)
}

type RequestGetUsecase interface {
	GetRequestsWithFilter(params RequestQueryParams) (*RequestListRes, error)
}