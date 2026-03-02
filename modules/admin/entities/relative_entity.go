package entities

import (

	"github.com/google/uuid"
)

type PatientDiseaseEntity struct {
	DiseaseID uuid.UUID `json:"disease_id"`
	Name      string    `json:"name"`
}

type RelativeEntity struct {
	Title       string     `json:"title"`
	FirstName   string     `json:"firstname"`
	LastName    string     `json:"lastname"`
	PhoneNumber string     `json:"phonenumber"`
	Address     Address `json:"address"`
	Role        string     `json:"role"`
	PatientID   uuid.UUID  `json:"patient_id"`
	CreatedBy   *uuid.UUID  `json:"created_by"`
	UpdatedBy   *uuid.UUID  `json:"updated_by"`
}

type RelativeCreate struct {
	Title       string     `json:"title"`
	FirstName   string     `json:"firstname"`
	LastName    string     `json:"lastname"`
	PhoneNumber string     `json:"phonenumber"`
	Address     Address `json:"address"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	UpdatedBy   uuid.UUID  `json:"updated_by"`
}

type RelativeInfo struct {
	Fullname    string `json:"fullname"`
	PhoneNumber string `json:"phonenumber"`
	Role        string `json:"role"`
	Address     Address `json:"address"`
}

type Relative struct {
	Kin       RelativeInfo `json:"kin"`
	Caretaker RelativeInfo `json:"caretaker"`
	Medicine  RelativeInfo `json:"medicine"`
}


type OfficerInfo struct {
	Fullname    string `json:"fullname"`
	Role        string `json:"role"`
}

type Officer struct {
	House OfficerInfo `json:"house"`
	Nurse OfficerInfo `json:"nurse"`
}

type RelativeAllUpdateReq struct {
	Kin       RelativeUpdateReq `json:"kin"`
	Caretaker RelativeUpdateReq `json:"caretaker"`
	Medicine  RelativeUpdateReq `json:"medicine"`
}

type RelativeUpdateReq struct {
	Title       string     `json:"title"`
	FirstName   string     `json:"firstname"`
	LastName    string     `json:"lastname"`
	PhoneNumber string     `json:"phonenumber"`
	Address     Address    `json:"address"`
}

type RelativeAllUpdateRes struct {
	Kin       RelativeUpdateRes `json:"kin"`
	Caretaker RelativeUpdateRes `json:"caretaker"`
	Medicine  RelativeUpdateRes `json:"medicine"`
}

type RelativeUpdateRes struct {
	ID          uuid.UUID `json:"id"`
	FullName    string     `json:"fullname"`
	PhoneNumber string     `json:"phonenumber"`
	Role        string     `json:"role"`
}

type OfficerAllUpdateReq struct {
	House     OfficerUpdateReq `json:"house"`
	Nurse     OfficerUpdateReq `json:"nurse"`
	UpdatedBy uuid.UUID        `json:"updated_by"`
}

type OfficerUpdateReq struct {
	Title     string  `json:"title"`
	FirstName string  `json:"firstname"`
	LastName  string  `json:"lastname"`
}

type OfficerAllUpdateRes struct {
	House OfficerUpdateRes `json:"house"`
	Nurse OfficerUpdateRes `json:"nurse"`
}

type OfficerUpdateRes struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"fullname"`
	Role     string    `json:"role"`
}

type RelativeCreateRepo interface {
	Create(relatives []RelativeEntity, adminID uuid.UUID) error
}

type RelativeUpdateRepo interface {
	UpdateRelativeInfo(patientID uuid.UUID, req *RelativeAllUpdateReq, adminID uuid.UUID) (*RelativeAllUpdateRes, error)
	UpdateOfficerInfo(patientID uuid.UUID, req *OfficerAllUpdateReq, adminID uuid.UUID) (*OfficerAllUpdateRes, error)
}

type RelativeUsecase interface {
	UpdateAllRelatives(patientID uuid.UUID, req *RelativeAllUpdateReq, adminID uuid.UUID) (*RelativeAllUpdateRes, error)
	UpdateAllOfficers(patientID uuid.UUID, req *OfficerAllUpdateReq, adminID uuid.UUID) (*OfficerAllUpdateRes, error)
}
