package usecases

import (
	"fmt"
	"time"
	"math"

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"
)

type AppointmentQueryUsecase struct {
	readRepo entities.AppointmentQueryRepo
}

func NewAppointmentQueryUsecase(readRepo entities.AppointmentQueryRepo) entities.AppointmentQueryUsecase {
	return &AppointmentQueryUsecase{readRepo: readRepo}
}

type appointmentCommandUsecase struct {
	repo entities.AppointmentCommandRepo
}

func NewAppointmentCommandUsecase(repo entities.AppointmentCommandRepo) entities.AppointmentCommandUsecase {
	return &appointmentCommandUsecase{repo: repo}
}

func (u *AppointmentQueryUsecase) GetAppointmentDetail(patientID uuid.UUID, diseaseID uuid.UUID) (*entities.AppointmentEntity, error) {
	hasDisease, err := u.readRepo.CheckPatientHasDisease(patientID, diseaseID)
	if err != nil {
		return nil, fmt.Errorf("error verifying patient medical records: %w", err)
	}
	if !hasDisease {
		return nil, fmt.Errorf("patient did not registered with this disease")
	}

	appointDB, adminDB, err := u.readRepo.GetAppointmentDetail(patientID, diseaseID)
	if err != nil {
		return nil, err
	}

	res := &entities.AppointmentEntity{
		ID:        appointDB.ID,
		No: 	   appointDB.No,	
		Doctor:    appointDB.Doctor,
		Status:    appointDB.Status,
		Purpose:   appointDB.Purpose,
		Place:     appointDB.Place,
		Date:      appointDB.Date.Format("2006-01-02"),
		StartTime: appointDB.StartTime,
		EndTime:   appointDB.EndTime,
		Symptom:   appointDB.Symptom,
		Note:      appointDB.Note,
		Delay:     appointDB.Delay, 
		DiseaseID: appointDB.DiseaseID,
		CreatedAt: appointDB.CreatedAt.Format("2006-01-02"),
	}

	if adminDB != nil && adminDB.FirstName != "" {
        res.CreatedBy = adminDB.Title + adminDB.FirstName + " " + adminDB.LastName
    }

	return res, nil
}

func (u *AppointmentQueryUsecase) ListPatientAppointments(patientID uuid.UUID) (*entities.PatientAppointment, error) {
	if patientID == uuid.Nil {
		return nil, fmt.Errorf("user ID cannot be nil")
	}

	appointDB, err := u.readRepo.ListPatientAppointments(patientID)
	if err != nil {
		return nil, err
	}

	if len(appointDB) == 0 {
		return nil, fmt.Errorf("patient has no appointments")
	}

	patient := appointDB[0].Patient
	if patient.ID == uuid.Nil {
		return nil, fmt.Errorf("patient data not loaded")
	}

	fullName := patient.Title + patient.FirstName + " " + patient.LastName

	var ageYears, ageMonths, ageDays int
	if !patient.DOB.IsZero() {
		now := time.Now()
		years := now.Year() - patient.DOB.Year()
		months := int(now.Month()) - int(patient.DOB.Month())
		days := now.Day() - patient.DOB.Day()

		if days < 0 {
			prevMonth := now.AddDate(0, -1, 0)
			days += utils.DaysInMonth(prevMonth.Year(), prevMonth.Month())
			months--
		}
		if months < 0 {
			months += 12
			years--
		}
		ageYears, ageMonths, ageDays = years, months, days
	}

	appointmentEntities := make([]entities.AppointmentEntity, len(appointDB))
	for i, app := range appointDB {
		appointmentEntities[i] = entities.AppointmentEntity{
			ID:          app.ID,
			No:          app.No,
			Doctor:      app.Doctor,
			Status:      app.Status,
			Purpose:     app.Purpose,
			Place:       app.Place,
			Date:        app.Date.Format("2006-01-02"),
			StartTime:   app.StartTime,
			EndTime:     app.EndTime,
			Symptom:     app.Symptom,
			Note:        app.Note,
			Delay:       app.Delay,
			DiseaseID:   app.DiseaseID,
			DiseaseName: app.Disease.Name,
		}
	}

	return &entities.PatientAppointment{
		PatientID:    patient.ID,
		FullName:     fullName,
		HnNumber:     patient.HnID,
		AgeYears:     ageYears,
		AgeMonths:    ageMonths,
		AgeDays:      ageDays,
		Appointments: appointmentEntities,
	}, nil

}

func (u *appointmentCommandUsecase) SubmitDelayRequest(userID uuid.UUID, patientID uuid.UUID, payload *entities.AppointmentDelayReq) error {
	exists, err := u.repo.CheckAppointmentExists(payload.AppointID, patientID)
	if err != nil {
		return fmt.Errorf("error verifying appointment: %w", err)
	}
	if !exists {
		return fmt.Errorf("appointment not found or does not belong to you")
	}

	parsedDate, err := time.Parse("2006-01-02", payload.Date)
	if err != nil {
		return fmt.Errorf("invalid date format, please use YYYY-MM-DD")
	}

	domainReq := &entities.DelayRequestEntity{
		RequestType: "delay",
		Description: "Patient requested to delay appointment",
		Date:        parsedDate,
		StartTime:   payload.StartTime,
		EndTime:     payload.EndTime,
		Status:      "pending",
		PatientID:   patientID,
		AppointID:   payload.AppointID,
		DiseaseID:   payload.DiseaseID,
		CreatedBy:   userID,
	}

	return u.repo.SaveDelayRequest(domainReq)
}

func (u *AppointmentQueryUsecase) GetScheduleByDisease(patientID uuid.UUID, diseaseID uuid.UUID) (*entities.DiseaseScheduleEntity, error) {
	if diseaseID == uuid.Nil {
		return nil, fmt.Errorf("disease ID is required")
	}

	hasDisease, err := u.readRepo.CheckPatientHasDisease(patientID, diseaseID)
	if err != nil {
		return nil, fmt.Errorf("error verifying patient medical records: %w", err)
	}
	if !hasDisease {
		return nil, fmt.Errorf("patient did not registered with this disease")
	}

	schedule, err := u.readRepo.GetScheduleByDisease(diseaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get calendar schedule: %w", err)
	}

	return schedule, nil
}

func (u *AppointmentQueryUsecase) GetDiseaseHistory(patientID uuid.UUID, diseaseID uuid.UUID, page int) (*entities.DiseaseHistoryResponse, error) {
	if page < 1 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	totalRows, err := u.readRepo.CountDiseaseHistory(patientID, diseaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to count history: %w", err)
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

	appointments, err := u.readRepo.GetDiseaseHistory(patientID, diseaseID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch history: %w", err)
	}

	return &entities.DiseaseHistoryResponse{
		TotalPages:   totalPages,
		CurrentPage:  page,
		Appointments: appointments, 
	}, nil
}

func (u *AppointmentQueryUsecase) GetHistoryDetail(appointID uuid.UUID, patientID uuid.UUID) (*entities.HistoryDetailEntity, error) {
	if appointID == uuid.Nil {
		return nil, fmt.Errorf("appointment ID is required")
	}

	detail, err := u.readRepo.GetHistoryDetail(appointID, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch appointment details: %w", err)
	}

	return detail, nil
}