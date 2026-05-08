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
		Preload("Disease").
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

func (r *AppointmentRepository) IsVaccineDisease(diseaseID uuid.UUID) (bool, error) {

	var disease databases.Disease

	err := r.db.
		Where("id = ?", diseaseID).
		First(&disease).Error
	if err != nil {
		return false, err
	}

	if disease.Name == "วัคซีน" {
		return true, nil
	}

	return false, nil
}

func (r *AppointmentRepository) CreateAppointment(e *entities.AppointmentEntity, adminID uuid.UUID) (*databases.Appoint, error) {
	var lastNo int
	err := r.db.Model(&databases.Appoint{}).
		Where("patient_id = ? AND disease_id = ?", e.PatientID, e.DiseaseID).
		Select("COALESCE(MAX(no), 0)").
		Scan(&lastNo).Error
	if err != nil {
		return nil, fmt.Errorf("query last appointment number: %w", err)
	}

	var parseDate time.Time
	if e.Date != "" {
		d, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			return nil, fmt.Errorf("parse date: %w", err)
		}
		parseDate = d
	}

	appoint := databases.Appoint{
		No:        lastNo + 1,
		DoctorID:  e.DoctorID,
		Status:    e.Status,
		Purpose:   e.Purpose,
		Place:     e.Place,
		StartTime: e.StartTime,
		EndTime:   e.EndTime,
		Date:      parseDate,
		PatientID: e.PatientID,
		DiseaseID: e.DiseaseID,
		Symptom:   e.Symptom,
		Note:      e.Note,
		Health: databases.Health{
			Weight:   e.Health.Weight,
			Height:   e.Health.Height,
			BMI:      e.Health.BMI,
			Pulse:    e.Health.Pulse,
			Sugar:    e.Health.Sugar,
			Pressure: e.Health.Pressure,
		},
		CreatedBy: &adminID,
		UpdatedBy: &adminID,
	}

	if err := r.db.Create(&appoint).Error; err != nil {
		return nil, fmt.Errorf("create appointment: %w", err)
	}

	r.UpdatePatientHealth(e.PatientID, e.Health.Weight, e.Health.Height, adminID)

	return &appoint, nil
}

func (r *AppointmentRepository) UpdateHealth(appointID uuid.UUID, health *entities.Health, adminID uuid.UUID) error {
    
    dbHealth := databases.Health{
        Weight:   health.Weight,
        Height:   health.Height,
        BMI:      health.BMI,
        Pulse:    health.Pulse,
        Sugar:    health.Sugar,
        Pressure: health.Pressure,
    }

    return r.db.Model(&databases.Appoint{}).
        Where("id = ?", appointID).
        Updates(map[string]interface{}{
            "health":     dbHealth,
            "updated_by": adminID,
        }).Error
}

func (r *AppointmentRepository) UpdatePatientHealth(patientID uuid.UUID, weight float64, height int, adminID uuid.UUID) error {

	return r.db.Model(&databases.Patient{}).
		Where("id = ?", patientID).
		Updates(map[string]interface{}{
			"weight":     weight,
			"height":     height,
			"updated_by": adminID,
		}).Error
}

func (r *AppointmentRepository) FindByID(appointID uuid.UUID) (*databases.Appoint, error) {
	var appoint databases.Appoint
	if err := r.db.
	Preload("Vaccinations").
	First(&appoint, "id = ?", appointID).Error; err != nil {
		return nil, err
	}
	return &appoint, nil
}

func (r *AppointmentRepository) UpdateSymptomNote(doctorID uuid.UUID, appointID uuid.UUID, symptom string, note string, adminID uuid.UUID) error {
	return r.db.Model(&databases.Appoint{}).
		Where("id = ?", appointID).
		Updates(map[string]interface{}{
			"symptom":   symptom,
			"note":      note,
			"doctor_id": doctorID,
			"updated_by": adminID,
		}).Error
}

func (r *AppointmentRepository) CompleteVaccinationRecord(recordID uuid.UUID, adminID uuid.UUID) error {
	return r.db.
		Model(&databases.VaccinationRecord{}).
		Where("id = ?", recordID).
		Updates(map[string]interface{}{
			"status":     "completed",
			"updated_by": adminID,
		}).Error
}

func (r *AppointmentRepository) UpdateVaccineDoctor(recordID uuid.UUID, doctorID uuid.UUID, adminID uuid.UUID) error {
	return r.db.Model(&databases.VaccinationRecord{}).
		Where("id = ?", recordID).
		Updates(map[string]interface{}{
			"vaccine_doctor_id": doctorID,
			"updated_by":        adminID,
		}).Error
}

func (r *AppointmentRepository) CreateVaccinationRecord(vaccineID uuid.UUID, appointID uuid.UUID, dose int, doctor uuid.UUID, adminID uuid.UUID) (error) {

	record := databases.VaccinationRecord{
		VaccineID:     vaccineID,
		AppointID:     appointID,
		DoseNumber:    dose,
		VaccineDoctor: doctor,
		Status:        "ongoing",
		CreatedBy:     &adminID,
		UpdatedBy:     &adminID,
	}

	return r.db.Create(&record).Error
}

func (r *AppointmentRepository) FindPatientVaccine(patientID uuid.UUID) (bool, error) {

	var count int64

	err := r.db.
		Table("patient_disease pd").
		Joins("JOIN disease d ON d.id = pd.disease_id").
		Where("pd.patient_id = ?", patientID).
		Where("d.name = ?", "วัคซีน").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AppointmentRepository) CheckVaccineExists(vaccineID uuid.UUID) (bool, error) {

	var count int64

	err := r.db.
		Model(&databases.Vaccine{}).
		Where("id = ?", vaccineID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AppointmentRepository) FindVaccineDiseaseID() (uuid.UUID, error) {

	var disease databases.Disease

	err := r.db.
		Where("name = ?", "วัคซีน").
		First(&disease).Error

	if err != nil {
		return uuid.Nil, err
	}

	return disease.ID, nil
}

func (r *AppointmentRepository) UpdateAppointment(appointID uuid.UUID, purpose string, place string, date string, startTime string, endTime string, doctorID uuid.UUID, adminID uuid.UUID) (*databases.Appoint, error) {
	
	var parsedDate time.Time
	if date != "" {
		parsedDate, _ = time.Parse("2006-01-02", date)
	}

	err := r.db.Model(&databases.Appoint{}).
		Where("id = ?", appointID).
		Updates(map[string]interface{}{
			"purpose":    purpose,
			"place":      place,
			"date":       parsedDate,
			"start_time": startTime,
			"end_time":   endTime,
			"doctor_id":  doctorID,
			"updated_by": adminID,
		}).Error

	if err != nil {
		return nil, err
	}

	var appoint databases.Appoint
	r.db.First(&appoint, "id = ?", appointID)
	return &appoint, nil
}

func (r *AppointmentRepository) FindMaxDose(patientID, vaccineID uuid.UUID) (int, error) {

	var dose int

	err := r.db.
		Model(&databases.VaccinationRecord{}).
		Joins("JOIN appoint ON appoint.id = vaccination_record.appoint_id").
		Where("appoint.patient_id = ? AND vaccination_record.vaccine_id = ?", patientID, vaccineID).
		Select("COALESCE(MAX(vaccination_record.dose_number),0)").
		Scan(&dose).Error

	return dose, err
}

func (r *AppointmentRepository) UpdateVaccinationRecord(appointID uuid.UUID, vaccineID uuid.UUID, dose int, adminID uuid.UUID) (error) {

	return r.db.Model(&databases.VaccinationRecord{}).
		Where("appoint_id = ?", appointID).
		Updates(map[string]interface{}{
			"vaccine_id":   vaccineID,
			"dose_number":  dose,
			"updated_by":   adminID,
		}).Error
}

func (r *AppointmentRepository) UpdateVaccineAppointment(appointID uuid.UUID, place string, date string, start string, end string, doctorID uuid.UUID, adminID uuid.UUID) (*databases.Appoint, error) {
	
	err := r.db.Model(&databases.Appoint{}).
		Where("id = ?", appointID).
		Updates(map[string]interface{}{
			"place":      place,
			"date":       date,
			"start_time": start,
			"end_time":   end,
			"doctor_id":  doctorID,
			"updated_by": adminID,
		}).Error

	if err != nil {
		return nil, err
	}

	var appoint databases.Appoint
	r.db.First(&appoint, "id = ?", appointID)
	return &appoint, nil
}

func (r *AppointmentRepository) FindAllDoctor() ([]databases.Doctor, error) {
	var doctors []databases.Doctor

	err := r.db.Select("id", "title", "first_name", "last_name").
		Order("first_name ASC").
		Find(&doctors).Error
	return doctors, err
}

