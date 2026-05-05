package repositories

import (
	"errors"
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
	"gorm.io/gorm"	
)

type requestCommandRepo struct {
	db *gorm.DB
}

func NewRequestCommandRepo(db *gorm.DB) entities.RequestCommandRepo {
	return &requestCommandRepo{db: db}
}

type requestQueryRepo struct {
	db *gorm.DB
}

func NewRequestQueryRepo(db *gorm.DB) entities.RequestQueryRepo {
	return &requestQueryRepo{db: db}
}

func (r *requestCommandRepo) GetLatestAppointID(patientID uuid.UUID, diseaseID uuid.UUID) (*uuid.UUID, error) {
	var app databases.Appoint
	err := r.db.Select("id").
		Where("patient_id = ? AND disease_id = ?", patientID, diseaseID).
		Order("date DESC, start_time DESC").
		First(&app).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no appointment history found for this disease")
		}
		return nil, err
	}
	return &app.ID, nil
}

func (r *requestCommandRepo) GetDiseaseName(diseaseID uuid.UUID) (string, error) {
	var disease databases.Disease
	err := r.db.Select("name").Where("id = ?", diseaseID).First(&disease).Error
	if err != nil {
		return "", err
	}
	return disease.Name, nil
}

func (r *requestCommandRepo) SaveMultipleRequests(reqs []entities.RequestEntity, notis []entities.NotificationEntity) error {

	tx := r.db.Begin()

	dbModels := make([]databases.Request, len(reqs))
	for i, req := range reqs {
		dbModels[i] = databases.Request{
			RequestType: req.RequestType,
			Description: req.Description,
			Date:        req.Date,
			Status:      req.Status,
			PatientID:   req.PatientID,
			CreatedBy:   &req.CreatedBy,
			AppointID:   req.AppointID,
			DiseaseID:   req.DiseaseID,
		}
	}
	if err := tx.Create(&dbModels).Error; err != nil {
			tx.Rollback()
			return err
		}

	dbNotis := make([]databases.Notification, len(notis))
	for i, noti := range notis {
		dbNotis[i] = databases.Notification{
			Header:    noti.Header,
			Body:      noti.Body,
			PatientID: noti.PatientID,
			CreatedBy: &noti.CreatedBy,
			
			RequestID: &dbModels[i].ID, 
		}
	}

	if len(dbNotis) > 0 {
		if err := tx.Create(&dbNotis).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *requestQueryRepo) CheckHasAnyAppointment(patientID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&databases.Appoint{}).
		Where("patient_id = ?", patientID).
		Count(&count).Error
	return count > 0, err
}

func (r *requestQueryRepo) CheckHasVaccineHistory(patientID uuid.UUID) (bool, error) {
	var count int64

	err := r.db.Table("vaccination_record").
		Joins("JOIN appoint ON appoint.id = vaccination_record.appoint_id").
		Where("appoint.patient_id = ?", patientID).
		Count(&count).Error
	return count > 0, err
}

func (r *requestQueryRepo) GetPendingRequests(patientID uuid.UUID) ([]entities.PendingRequestData, error) {
	var pending []entities.PendingRequestData

	err := r.db.Model(&databases.Request{}).
		Select("request_type, description, disease_id").
		Where("patient_id = ? AND status = ?", patientID, "pending").
		Find(&pending).Error

	return pending, err
}

func (r *requestQueryRepo) GetLatestCompletedAppoint(patientID uuid.UUID) (*uuid.UUID, *uuid.UUID, error) {
	var app databases.Appoint

	err := r.db.Select("id", "disease_id").
		Where("patient_id = ? AND status = ?", patientID, "completed"). 
		Order("date DESC, start_time DESC").
		Limit(1).
		Find(&app).Error

	if err != nil {
		return nil, nil, err
	}

	if app.ID == uuid.Nil {
		return nil, nil, nil 
	}

	return &app.ID, &app.DiseaseID, nil
}

