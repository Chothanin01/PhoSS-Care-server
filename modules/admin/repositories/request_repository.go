package repositories

import (
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
	"gorm.io/gorm"
)

type RequestGetRepository struct {
	db *gorm.DB
}

func NewRequestGetRepository(db *gorm.DB) *RequestGetRepository {
	return &RequestGetRepository{db: db}
}

func (r *RequestGetRepository) GetRequestsWithFilter(params entities.RequestQueryParams) ([]entities.RequestInfo, error) {
	var results []entities.RequestInfo

	query := r.db.
		Table("request AS r").
		Select(`
			r.id AS id,
			r.request_type AS request_type,
			r.status AS status,
			p.hn_id AS hn_number,
			CONCAT(p.title, p.first_name, ' ', p.last_name) AS patient_name,
			COALESCE(d.name, '') AS disease_name,
			r.description AS description,
			r.date AS date,
			r.time AS time
		`).
		Joins("JOIN patient p ON p.id = r.patient_id").
		Joins("LEFT JOIN appoint a ON a.id = r.appoint_id").
		Joins("LEFT JOIN disease d ON d.id = a.disease_id").
		Order("r.created_at DESC")

	if params.ReqType != "" {
		query = query.Where("r.request_type ILIKE ?", "%"+params.ReqType+"%")
	}

	if params.Status != "" {
		query = query.Where("r.status ILIKE ?", "%"+params.Status+"%")
	}

	offset := (params.Page - 1) * params.Limit
	if offset < 0 {
		offset = 0
	}

	err := query.
		Order("r.created_at DESC").
		Limit(params.Limit).
		Offset(offset).
		Scan(&results).Error

	return results, err
}

func (r *RequestGetRepository) CountRequestsWithFilter(params entities.RequestQueryParams) (int64, error) {
	var count int64
	query := r.db.Model(&databases.Request{})

	if params.ReqType != "" {
		query = query.Where("request_type ILIKE ?", "%"+params.ReqType+"%")
	}

	if params.Status != "" {
		query = query.Where("status ILIKE ?", "%"+params.Status+"%")
	}

	err := query.Count(&count).Error
	return count, err
}