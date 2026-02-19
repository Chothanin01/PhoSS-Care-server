package usecases

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
)

type RelativeUpdateUsecase interface {
	UpdateAllRelatives(patientID uuid.UUID, req *entities.RelativeAllUpdateReq) (*entities.RelativeAllUpdateRes, error)
}

type relativeUpdateUsecase struct {
	repo entities.RelativeUpdateRepo
}

func NewRelativeUsecase(repo entities.RelativeUpdateRepo) entities.RelativeUsecase {
	return &relativeUpdateUsecase{repo: repo}
}

func (u *relativeUpdateUsecase) UpdateAllRelatives(patientID uuid.UUID, req *entities.RelativeAllUpdateReq) (*entities.RelativeAllUpdateRes, error) {
	if req.UpdatedBy == uuid.Nil {
		return nil, fmt.Errorf("missing updated_by field")
	}
	return u.repo.UpdateRelativeInfo(patientID, req)
}