package usecases

import (
	"fmt"
	"math"

	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
	"github.com/google/uuid"
)

type vaccineQueryUsecase struct {
	repo entities.VaccineQueryRepo
}

func NewVaccineQueryUsecase(repo entities.VaccineQueryRepo) entities.VaccineQueryUsecase {
	return &vaccineQueryUsecase{repo: repo}
}

func (u *vaccineQueryUsecase) GetVaccineList(patientID uuid.UUID, page int, filter string) (*entities.VaccineListResponse, error) {
	if patientID == uuid.Nil {
		return nil, fmt.Errorf("patient ID is required")
	}

	if page < 1 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	validFilters := map[string]bool{"all": true, "completed": true, "ongoing": true, "not_vaccinated": true}
	if !validFilters[filter] {
		filter = "all" 
	}

	totalRows, err := u.repo.CountVaccines(patientID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count vaccines: %w", err)
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

	vaccines, err := u.repo.GetVaccineList(patientID, filter, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vaccine list: %w", err)
	}

	if vaccines == nil {
		vaccines = make([]entities.VaccineItemEntity, 0)
	}

	return &entities.VaccineListResponse{
		TotalPages:  totalPages,
		CurrentPage: page,
		Vaccines:    vaccines,
	}, nil
}