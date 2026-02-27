package usecases

import (
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
)

type requestGetUsecase struct {
	repo entities.RequestGetRepo
}

func NewRequestGetUsecase(repo entities.RequestGetRepo) entities.RequestGetUsecase {
	return &requestGetUsecase{repo: repo}
}

func (u *requestGetUsecase) GetRequestsWithFilter(params entities.RequestQueryParams) (*entities.RequestListRes, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	data, err := u.repo.GetRequestsWithFilter(params)
	if err != nil {
		return nil, err
	}

	total, err := u.repo.CountRequestsWithFilter(params)
	if err != nil {
		return nil, err
	}

	return &entities.RequestListRes{
		Page:       params.Page,
		PerPage:    params.Limit,
		TotalPages: int((total + int64(params.Limit) - 1) / int64(params.Limit)),
		Data:       data,
	}, nil
}