package repositories

import (
	"errors"
	"time"
    "fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
)

type AppointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) FindOngoing(patientID, diseaseID uuid.UUID) (*databases.Appoint, error) {
	var appoint databases.Appoint
	err := r.db.
		Where("patient_id = ? AND disease_id = ? AND status = ?", patientID, diseaseID, "ongoing").
		Order("created_at DESC").
		First(&appoint).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &appoint, err
}

func (r *AppointmentRepository) FindOngoingVaccination(patientID uuid.UUID) (*databases.VaccinationRecord, error) {
	var record databases.VaccinationRecord
    
	err := r.db.
		Joins("JOIN appoint ON appoint.id = vaccination_record.appoint_id").
		Where("appoint.patient_id = ?", patientID).
		Where("vaccination_record.status = ?", "ongoing").
		Preload("Vaccine").
		Order("vaccination_record.dose_number DESC").
		First(&record).Error

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *AppointmentRepository) CompleteAppoint(appointID, adminID uuid.UUID) error {
    return r.db.Model(&databases.Appoint{}).
        Where("id = ? AND status = ?", appointID, "ongoing").
        Updates(map[string]interface{}{
            "status":     "completed",
            "updated_by": adminID,
        }).Error
}

func (r *AppointmentRepository) DiseaseExists(diseaseID uuid.UUID) (bool, error) {
    var count int64
    if err := r.db.Model(&databases.Disease{}).
        Where("id = ?", diseaseID).
        Count(&count).Error; err != nil {
        return false, err
    }
    return count > 0, nil
}

func (r *AppointmentRepository) CreateAppointment(e *entities.AppointmentEntity, adminID uuid.UUID) (*databases.Appoint, error) {
    var lastNo int
    err := r.db.
        Model(&databases.Appoint{}).
        Where("patient_id = ? AND disease_id = ?", e.PatientID, e.DiseaseID).
        Select("COALESCE(MAX(no), 0)").
        Scan(&lastNo).Error
    if err != nil {
        return nil, fmt.Errorf("query last appointment number: %w", err)
    }

    No := lastNo + 1

    var parseDate time.Time
	if e.Date != "" {
		d, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			return nil, fmt.Errorf("parse date: %w", err)
		}
		parseDate = d
	}

    appoint := databases.Appoint{
        No:        No,
        Doctor:    e.Doctor,
        Status:    e.Status,
        Purpose:   e.Purpose,
        Place:     e.Place,
        Time:      e.Time,
        Date:      parseDate,
        PatientID: e.PatientID,
        DiseaseID: e.DiseaseID,
		Symptom: e.Symtom,
		Note: e.Note,
        CreatedBy: &adminID,
        UpdatedBy: &adminID,
    }

    if err := r.db.Create(&appoint).Error; err != nil {
        return nil, fmt.Errorf("create appointment: %w", err)
    }

    return &appoint, nil
}

func (r *AppointmentRepository) CreateHealthRecord(health *entities.Health, patientID uuid.UUID, appointID uuid.UUID, adminID uuid.UUID) error {
    healthRecord := databases.Health{
        Weight:    health.Weight,
        Height:    health.Height,
        BMI:       health.BMI,
        Pulse:     health.Pulse,
        Sugar:     health.Sugar,
        PatientID: patientID,
        AppointID: appointID,
        CreatedBy: &adminID,
        UpdatedBy: &adminID,
    }
	return r.db.Create(&healthRecord).Error
}

func (r *AppointmentRepository) FindByID(appointID uuid.UUID) (*databases.Appoint, error) {
	var appoint databases.Appoint
	if err := r.db.Preload("Healths").First(&appoint, "id = ?", appointID).Error; err != nil {
		return nil, err
	}
	return &appoint, nil
}

func (r *AppointmentRepository) UpdateSymptomNote(doctor string,appointID uuid.UUID, symptom string, note string, adminID uuid.UUID) error {
	return r.db.
		Model(&databases.Appoint{}).
		Where("id = ?", appointID).
		Updates(map[string]interface{}{
			"symptom":    symptom,
			"note":       note,
			"doctor":     doctor,
			"updated_by": adminID,
		}).Error
}

func (r *AppointmentRepository) UpdateHealth(appointID uuid.UUID, health *entities.Health, adminID uuid.UUID) error {
	var existing databases.Health
	err := r.db.Where("appoint_id = ?", appointID).First(&existing).Error

	if err == nil {
		return r.db.Model(&existing).Updates(map[string]interface{}{
			"weight":     health.Weight,
			"height":     health.Height,
			"bmi":        health.BMI,
			"pulse":      health.Pulse,
			"sugar":      health.Sugar,
			"updated_by": adminID,
		}).Error
	}
	return err
}

