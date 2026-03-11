package repositories

import (

	"github.com/google/uuid"
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

type RequestUpdateRepository struct {
	db *gorm.DB
}

func NewRequestUpdateRepository(db *gorm.DB) *RequestUpdateRepository {
	return &RequestUpdateRepository{db: db}
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
			r.time AS time,
			r.appoint_id AS appoint_id
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

func (r *RequestGetRepository) GetRequestInfoByID(id uuid.UUID) (*databases.Request, error) {
	var req databases.Request

	err := r.db.
		Preload("Patient").
		Preload("Appoint").
		Preload("Appoint.Disease").
		First(&req, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &req, nil
}
func (r *RequestUpdateRepository) FindRequestByID(id uuid.UUID) (*entities.RequestInfo, error) {
	var req databases.Request
	err := r.db.Preload("Appoint").First(&req, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	detail := &entities.RequestInfo{
		ID:          req.ID,
		RequestType: req.RequestType,
		Status:      req.Status,
		Description: req.Description,
	}
	if req.AppointID != uuid.Nil {
		detail.AppointID = req.AppointID
	}
	return detail, nil
}

func (r *RequestUpdateRepository) UpdateRequestStatus(id uuid.UUID, status, description string, adminID uuid.UUID) error {
	return r.db.Model(&databases.Request{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"description": description,
			"updated_by":  adminID,
		}).Error
}

func (r *RequestUpdateRepository) UpdateAppointForAccepted(appointID uuid.UUID, adminID uuid.UUID) error {
	return r.db.Model(&databases.Appoint{}).
		Where("id = ?", appointID).
		Updates(map[string]interface{}{
			"delay":      true,
			"status":     "ongoing",
			"updated_by": adminID,
		}).Error
}