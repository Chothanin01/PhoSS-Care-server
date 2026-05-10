package repositories

import (
	"errors"
	"time"

	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type appointmentQueryRepo struct {
	db *gorm.DB
}

func NewappointmentQueryRepo(db *gorm.DB) *appointmentQueryRepo {
	return &appointmentQueryRepo{db: db}
}

type appointmentCommandRepo struct {
	db *gorm.DB
}

func NewAppointmentCommandRepo(db *gorm.DB) *appointmentCommandRepo {
	return &appointmentCommandRepo{db: db}
}

func (r *appointmentQueryRepo) GetAppointmentDetail(patientID uuid.UUID, diseaseID uuid.UUID) (*entities.AppointmentEntity, error) {
	var dbAppoint databases.Appoint

	err := r.db.
		Preload("Disease").
		Preload("Doctor").
		Where("patient_id = ? AND disease_id = ? AND status IN ?", patientID, diseaseID, []string{"ongoing", "delay"}).
		First(&dbAppoint).Error
		
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}

	detail := &entities.AppointmentEntity{
		ID:          dbAppoint.ID,
		No:          dbAppoint.No,
		Doctor:      dbAppoint.Doctor.Title + dbAppoint.Doctor.FirstName + " " + dbAppoint.Doctor.LastName,
		Status:      dbAppoint.Status,
		Purpose:     dbAppoint.Purpose,
		Place:       dbAppoint.Place,
		Date:        dbAppoint.Date.Format("2006-01-02"), 
		StartTime:   dbAppoint.StartTime,
		EndTime:     dbAppoint.EndTime,
		Symptom:     dbAppoint.Symptom,
		Note:        dbAppoint.Note,
		DiseaseID:   dbAppoint.DiseaseID,
		DiseaseName: dbAppoint.Disease.Name,
		CreatedAt:   dbAppoint.CreatedAt.Format(time.RFC3339),
	}

	if dbAppoint.CreatedBy != nil {
		var dbAdmin databases.Admin
		err = r.db.Select("title, first_name, last_name").Where("user_id = ?", *dbAppoint.CreatedBy).First(&dbAdmin).Error
		if err == nil {
			detail.CreatedBy = dbAdmin.Title + dbAdmin.FirstName + " " + dbAdmin.LastName
		} else {
			detail.CreatedBy = "Unknown Admin" 
		}
	}

	if dbAppoint.Status == "delay" {
		
		var req databases.Request
		err := r.db.Where("appoint_id = ? AND status = 'pending'", dbAppoint.ID).First(&req).Error
		
		if err == nil {
			detail.DelayDate = req.Date.Format("2006-01-02") 
			detail.DelayStartTime = req.StartTime
			detail.DelayEndTime = req.EndTime
		} 

	}

	return detail, nil
}

func (r *appointmentQueryRepo) GetPatientBasicInfo(patientID uuid.UUID) (*entities.PatientBasicInfo, error) {
	var dbPatient databases.Patient 

	err := r.db.Select("id, title, first_name, last_name, hn_id, dob").
		Where("id = ?", patientID).
		First(&dbPatient).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrNotFound 
		}
		return nil, err
	}
	cleanPatient := &entities.PatientBasicInfo{
		ID:        dbPatient.ID,
		Title:     dbPatient.Title,
		FirstName: dbPatient.FirstName,
		LastName:  dbPatient.LastName,
		HnID:      dbPatient.HnID,
		DOB:       dbPatient.DOB,
	}

	return cleanPatient, nil
}


func (r *appointmentQueryRepo) ListPatientAppointments(patientID uuid.UUID) ([]entities.AppointmentEntity, error) {
	var dbAppoints []databases.Appoint

	err := r.db.
		Preload("Disease").
		Preload("Patient").
		Preload("Doctor").
		Where("patient_id = ? AND status IN ?", patientID, []string{"ongoing", "delay"}).
		Find(&dbAppoints).
		Error

	if err != nil {
		return nil, err
	}

	var resultList []entities.AppointmentEntity

	for _, appt := range dbAppoints {
		
		entity := entities.AppointmentEntity{
			ID:          appt.ID,
			No:          appt.No,
			Doctor:      appt.Doctor.Title + appt.Doctor.FirstName + " " + appt.Doctor.LastName,
			Status:      appt.Status,
			Purpose:     appt.Purpose,
			Place:       appt.Place,
			Date:        appt.Date.Format("2006-01-02"),
			StartTime:   appt.StartTime,
			EndTime:     appt.EndTime,
			Symptom:     appt.Symptom,
			Note:        appt.Note,
			DiseaseID:   appt.DiseaseID,
			DiseaseName: appt.Disease.Name,
		}

		if appt.Status == "delay" {
			
			var req databases.Request
			err := r.db.Where("appoint_id = ? AND status = 'pending'", appt.ID).First(&req).Error
			
			if err != nil {
			} else {
				entity.DelayDate = req.Date.Format("2006-01-02")
				entity.DelayStartTime = req.StartTime
				entity.DelayEndTime = req.EndTime
			}
		}

		resultList = append(resultList, entity)
	}

	return resultList, nil
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

func (r *appointmentCommandRepo) GetOngoingAppointmentIDByDisease(patientID uuid.UUID, diseaseID uuid.UUID) (*uuid.UUID, error) {
	
	var result struct {
		ID uuid.UUID
	}
	
	err := r.db.Model(&databases.Appoint{}).
		Select("id").
		Where("patient_id = ? AND disease_id = ? AND status = ?", patientID, diseaseID, "ongoing").
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil 
		}
		return nil, err 
	}
	
	return &result.ID, nil
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

	return r.db.Transaction(func(tx *gorm.DB) error {
		
		if err := tx.Create(dbModel).Error; err != nil {
			return err
		}

		err := tx.Model(&databases.Appoint{}).
			Where("id = ? AND patient_id = ?", req.AppointID, req.PatientID).
			Updates(map[string]interface{}{
				"status":     "delay",
				"updated_at": time.Now(),
			}).Error

		if err != nil {
			return err
		}

		return nil
	})
}

func (r *appointmentQueryRepo) CheckPatientHasDisease(patientID uuid.UUID, diseaseID uuid.UUID) (bool, error) {
	var count int64
	
	err := r.db.Table("patient_disease").
		Where("patient_id = ? AND disease_id = ?", patientID, diseaseID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	
	return count > 0, nil 
}

func (r *appointmentQueryRepo) GetScheduleByDisease(patientID uuid.UUID, diseaseID uuid.UUID) (*entities.DiseaseScheduleEntity, error) {
	var dbDisease databases.Disease

	err := r.db.Where("id = ?", diseaseID).First(&dbDisease).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("disease not found")
		}
		return nil, err
	}

	var appointDate *time.Time
	err = r.db.Model(&databases.Appoint{}).
		Select("date").
		Where("patient_id = ? AND disease_id = ? AND status = ?", patientID, diseaseID, "ongoing").
		First(&appointDate).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	currentDateStr := ""
	if appointDate != nil {
		currentDateStr = appointDate.Format("2006-01-02")
	}

	domainSchedule := &entities.DiseaseScheduleEntity{
		DiseaseID:     dbDisease.ID,
		DiseaseName:   dbDisease.Name,
		AvailableDays: dbDisease.AvailableDays, 
		CurrentDate:   currentDateStr,
	}

	return domainSchedule, nil
}

func (r *appointmentQueryRepo) CountDiseaseHistory(patientID uuid.UUID, diseaseID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&databases.Appoint{}).
		Where("patient_id = ? AND disease_id = ?", patientID, diseaseID).
		Count(&count).Error
	return count, err
}

func (r *appointmentQueryRepo) GetDiseaseHistory(patientID uuid.UUID, diseaseID uuid.UUID, limit int, offset int) ([]entities.HistoryAppointEntity, error) {
	var dbAppoints []databases.Appoint

	err := r.db.
		Preload("CreatedByUser.Admin").
		Where("patient_id = ? AND disease_id = ?", patientID, diseaseID).
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&dbAppoints).Error

	if err != nil {
		return nil, err
	}

	domainAppoints := make([]entities.HistoryAppointEntity, len(dbAppoints))

	for i, app := range dbAppoints {
		doctorName := "Unknown"
		if app.Doctor.FirstName != "" {
			doctorName = app.Doctor.Title + app.Doctor.FirstName + " " + app.Doctor.LastName
		} else if app.CreatedByUser != nil && app.CreatedByUser.Admin != nil && app.CreatedByUser.Admin.FirstName != "" {
			admin := app.CreatedByUser.Admin
			doctorName = admin.Title + admin.FirstName + " " + admin.LastName
		}

		color := app.ColorStatus
		if color == "" {
			color = "none" 
		}

		domainAppoints[i] = entities.HistoryAppointEntity{
			AppointID:   app.ID,
			No:          app.No,
			Date:        app.Date.Format("2006-01-02"),
			Note:        app.Note,
			ColorStatus: color,
			DoctorName:  doctorName,
		}
	}

	return domainAppoints, nil
}

func (r *appointmentQueryRepo) GetHistoryDetail(appointID uuid.UUID, patientID uuid.UUID) (*entities.HistoryDetailEntity, error) {
	var app databases.Appoint

	err := r.db.
		Preload("CreatedByUser.Admin").
		Where("id = ? AND patient_id = ?", appointID, patientID).
		First(&app).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("appointment not found or access denied")
		}
		return nil, err
	}

	var prevApps []databases.Appoint
	var prevID *uuid.UUID
	r.db.Select("id").
		Where("patient_id = ? AND disease_id = ? AND (date < ? OR (date = ? AND start_time < ?))", 
			patientID, app.DiseaseID, app.Date, app.Date, app.StartTime).
		Order("date DESC, start_time DESC"). 
		Limit(1).
		Find(&prevApps)

	if len(prevApps) > 0 {
		prevID = &prevApps[0].ID
	}

	var nextApps []databases.Appoint
	var nextID *uuid.UUID
	r.db.Select("id").
		Where("patient_id = ? AND disease_id = ? AND (date > ? OR (date = ? AND start_time > ?))", 
			patientID, app.DiseaseID, app.Date, app.Date, app.StartTime).
		Order("date ASC, start_time ASC").
		Limit(1).
		Find(&nextApps)

	if len(nextApps) > 0 {
		nextID = &nextApps[0].ID
	}

	doctorName := "Unknown Doctor"
	if app.Doctor.FirstName != "" {
		doctorName = app.Doctor.Title + app.Doctor.FirstName + " " + app.Doctor.LastName
	} else if app.CreatedByUser != nil && app.CreatedByUser.Admin != nil && app.CreatedByUser.Admin.FirstName != "" {
		admin := app.CreatedByUser.Admin
		doctorName = admin.Title + admin.FirstName + " " + admin.LastName
	}  

	color := app.ColorStatus
	if color == "" {
		color = "none"
	}

	detail := &entities.HistoryDetailEntity{
		AppointID:   app.ID,
		No:          app.No,
		Date:        app.Date.Format("2006-01-02"),
		Note:        app.Note,
		ColorStatus: color,
		Doctor:      doctorName,
		Purpose:     app.Purpose,
		Symptom:     app.Symptom,
		Health: entities.HealthEntity{
			Pulse:    app.Health.Pulse,
			Pressure: app.Health.Pressure,
			Height:   app.Health.Height,
			Weight:   app.Health.Weight,
			BMI:      app.Health.BMI,
			Sugar:    app.Health.Sugar,
		},
		NextAppointID: nextID,
		PrevAppointID: prevID,
	}

	return detail, nil
}

func (r *appointmentCommandRepo) CancelDelayRequest(appointID uuid.UUID, patientID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		
		result := tx.Model(&databases.Request{}).
			Where("appoint_id = ? AND patient_id = ? AND request_type = ? AND status = ?", 
				appointID, patientID, "appoint", "pending").
			Update("status", "canceled")

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return entities.ErrNotFound
		}

		err := tx.Model(&databases.Appoint{}).
			Where("id = ? AND patient_id = ?", appointID, patientID).
			Update("status", "ongoing").Error

		if err != nil {
			return err
		}

		return nil
	})
}