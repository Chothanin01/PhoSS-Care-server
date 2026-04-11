package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
	"gorm.io/gorm"	
)

type appointmentQueryRepository struct {
	db *gorm.DB
}

func NewAppointmentQueryRepository(db *gorm.DB) *appointmentQueryRepository {
	return &appointmentQueryRepository{db: db}
}

type appointmentCommandRepo struct {
	db *gorm.DB
}

func NewAppointmentCommandRepo(db *gorm.DB) *appointmentCommandRepo {
	return &appointmentCommandRepo{db: db}
}

func (r *appointmentQueryRepository) GetAppointmentDetail(patientID uuid.UUID, diseaseID uuid.UUID) (*databases.Appoint, *databases.Admin, error) {
	var appoint databases.Appoint
	var admin databases.Admin

	err := r.db.
		Preload("Disease").
        Where("patient_id = ? AND disease_id = ? AND status IN ?", patientID, diseaseID, []string{"ongoing", "delay"}).
        First(&appoint).Error
    if err != nil {
        return nil, nil, err
    }
	
	if appoint.CreatedBy != nil {
        err = r.db.Where("user_id = ?", *appoint.CreatedBy).First(&admin).Error
        if err != nil {
            return &appoint, nil, nil 
        }
    }

	return &appoint, &admin, nil
}

func (r *appointmentQueryRepository) ListPatientAppointments(patientID uuid.UUID) ([]databases.Appoint, error) {
	var appoint []databases.Appoint

	err := r.db.
		Preload("Disease").
		Preload("Patient").
    	Where("patient_id = ? AND status IN ?", patientID, []string{"ongoing", "dalay"}).
    	Find(&appoint).
		Error

	if err != nil {
		return nil,err
	}

	return appoint, nil
}

func (r *appointmentCommandRepo) CheckAppointmentExists(appointID uuid.UUID, patientID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&databases.Appoint{}).
		Where("id = ? AND patient_id = ?", appointID, patientID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *appointmentCommandRepo) SaveDelayRequest(req *entities.DelayRequestEntity) error {
	dbModel := &databases.Request{
		RequestType: req.RequestType,
		Description: req.Description,
		Date:        req.Date,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Status:      req.Status,
		PatientID:   req.PatientID,
		AppointID:   req.AppointID,
		DiseaseID:   req.DiseaseID,
		CreatedBy:   &req.CreatedBy,
	}

	return r.db.Create(dbModel).Error
}

func (r *appointmentQueryRepository) CheckPatientHasDisease(patientID uuid.UUID, diseaseID uuid.UUID) (bool, error) {
	var count int64
	
	err := r.db.Table("patient_disease").
		Where("patient_id = ? AND disease_id = ?", patientID, diseaseID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	
	return count > 0, nil 
}

func (r *appointmentQueryRepository) GetScheduleByDisease(diseaseID uuid.UUID) (*entities.DiseaseScheduleEntity, error) {
	var dbDisease databases.Disease

	err := r.db.Where("id = ?", diseaseID).First(&dbDisease).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("disease not found")
		}
		return nil, err
	}

	domainSchedule := &entities.DiseaseScheduleEntity{
		DiseaseID:     dbDisease.ID,
		DiseaseName:   dbDisease.Name,
		AvailableDays: dbDisease.AvailableDays, 
	}

	return domainSchedule, nil
}