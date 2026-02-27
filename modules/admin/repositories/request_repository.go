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

func (r *RequestGetRepository) GetRequestInfoByID(id uuid.UUID) (*entities.RequestInfoRes, error) {
	var req databases.Request
	err := r.db.
		Preload("Appoint.Patient").
		Preload("Appoint").
		Preload("Patient").
		First(&req, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	p := req.Patient
	fullname := p.Title + p.FirstName + " " + p.LastName

	res := &entities.RequestInfoRes{
		RequestID: req.ID,
		RequestType: req.RequestType,
		Status:      req.Status,
		PatientID: 	 p.ID,
		FullName: fullname,
		IDCard:      p.IDCard,
		HnNumber:    p.HnID,
	}

	switch req.RequestType {
	case "appoint":
		res.Date = req.Date.Format("2006-01-02")
		res.Time = req.Time
		res.Doctor = req.Appoint.Doctor
		res.AppointDate = req.Appoint.Date.Format("2006-01-02")
		res.AppointTime = req.Appoint.Time

	case "medical":
		res.CreatedDate = req.CreatedAt.Format("2006-01-02")
		res.Description = req.Description
		res.Doctor = req.Appoint.Doctor

	case "document":
		res.CreatedDate = req.CreatedAt.Format("2006-01-02")
		res.Description = req.Description
	}

	return res, nil
}