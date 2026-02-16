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
	CreatedBy   uuid.UUID  `json:"created_by"`
	UpdatedBy   uuid.UUID  `json:"updated_by"`
}

type RelativeDetail struct {
	Title       string     `json:"title"`
	FirstName   string     `json:"firstname"`
	LastName    string     `json:"lastname"`
	PhoneNumber string     `json:"phonenumber"`
	Address     Address `json:"address"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	UpdatedBy   uuid.UUID  `json:"updated_by"`
}

type RelativeCreateRepo interface {
	Create(relatives []RelativeEntity) error
}
