package entities

import (
	"github.com/google/uuid"
)

type VaccineItemEntity struct {
	VaccineID        uuid.UUID `json:"vaccine_id"` 
	Name             string `json:"name"`
	Type             string `json:"type"`
	Age              string `json:"age"`
	VaccinatedDate   string `json:"vaccinated_date"`   
	VaccinatedStatus string `json:"vaccinated_status"` 
}

type VaccineListResponse struct {
	TotalPages  int                 `json:"total_pages"`
	CurrentPage int                 `json:"current_page"`
	Vaccines    []VaccineItemEntity `json:"vaccines"`
}

type VaccineDetailEntity struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	Age              string    `json:"age"`
	Effect           string    `json:"effect"` 
	Note             string    `json:"note"`   
	
	VaccinatedDate   string    `json:"vaccinated_date"`
	VaccinatedStatus string    `json:"vaccinated_status"`
	DoseNumber       int       `json:"dose_number,omitempty"`
	VaccineDoctor    string    `json:"vaccine_doctor,omitempty"`
}


type VaccineQueryRepo interface {
	CountVaccines(patientID uuid.UUID, filter string) (int64, error)
	GetVaccineList(patientID uuid.UUID, filter string, limit int, offset int) ([]VaccineItemEntity, error)

	GetVaccineDetail(patientID uuid.UUID, vaccineID uuid.UUID) (*VaccineDetailEntity, error)
}

type VaccineQueryUsecase interface {
	GetVaccineList(patientID uuid.UUID, page int, filter string) (*VaccineListResponse, error)

	GetVaccineDetail(patientID uuid.UUID, vaccineID uuid.UUID) (*VaccineDetailEntity, error)
}