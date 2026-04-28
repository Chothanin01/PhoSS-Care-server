package usecases

import (
	"fmt"
	"time"
	"math"

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"
)

type appointmentQueryUsecase struct {
	readRepo entities.AppointmentQueryRepo
}

func NewAppointmentQueryUsecase(readRepo entities.AppointmentQueryRepo) entities.AppointmentQueryUsecase {
	return &appointmentQueryUsecase{readRepo: readRepo}
}

type appointmentCommandUsecase struct {
	repo entities.AppointmentCommandRepo
}

func NewAppointmentCommandUsecase(repo entities.AppointmentCommandRepo) entities.AppointmentCommandUsecase {
	return &appointmentCommandUsecase{repo: repo}
}

func (u *appointmentQueryUsecase) GetAppointmentDetail(patientID uuid.UUID, diseaseID uuid.UUID) (*entities.AppointmentEntity, error) {
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



func (u *appointmentQueryUsecase) ListPatientAppointments(patientID uuid.UUID) (*entities.PatientAppointment, error) {
	
	if patientID == uuid.Nil {
		return nil, fmt.Errorf("user ID cannot be nil")
	}

	patient, err := u.readRepo.GetPatientBasicInfo(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to load patient data: %w", err)
	}

	appointments, err := u.readRepo.ListPatientAppointments(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to load appointments: %w", err)
	}

	if appointments == nil {
		appointments = make([]entities.AppointmentEntity, 0)
	}

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

	return &entities.PatientAppointment{
		PatientID:    patient.ID,
		FullName:     patient.Title + patient.FirstName + " " + patient.LastName,
		HnNumber:     patient.HnID,
		AgeYears:     ageYears,
		AgeMonths:    ageMonths,
		AgeDays:      ageDays,
		Appointments: appointments, 
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
		RequestType: "appoint",
		Description: "Patient requested to delay appointment",
		Date:        parsedDate,
		StartTime:   payload.StartTime,
		EndTime:     payload.EndTime,
		Status:      "pending",
		PatientID:   patientID,
		AppointID:   &payload.AppointID,
		DiseaseID:   &payload.DiseaseID,
		CreatedBy:   userID,
	}

	return u.repo.SaveDelayRequest(domainReq)
}

func (u *appointmentQueryUsecase) GetScheduleByDisease(patientID uuid.UUID, diseaseID uuid.UUID) (*entities.DiseaseScheduleEntity, error) {
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

	schedule, err := u.readRepo.GetScheduleByDisease(patientID, diseaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get calendar schedule: %w", err)
	}

	return schedule, nil
}

func (u *appointmentQueryUsecase) GetDiseaseHistory(patientID uuid.UUID, diseaseID uuid.UUID, page int) (*entities.DiseaseHistoryResponse, error) {
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

func (u *appointmentQueryUsecase) GetHistoryDetail(appointID uuid.UUID, patientID uuid.UUID) (*entities.HistoryDetailEntity, error) {
	if appointID == uuid.Nil {
		return nil, fmt.Errorf("appointment ID is required")
	}

	detail, err := u.readRepo.GetHistoryDetail(appointID, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch appointment details: %w", err)
	}

	return detail, nil
}

func (u *appointmentCommandUsecase) CancelDelayRequest(appointID uuid.UUID, patientID uuid.UUID) error {
	if appointID == uuid.Nil {
		return fmt.Errorf("appointment ID is required")
	}
	if patientID == uuid.Nil {
		return fmt.Errorf("patient ID is required")
	}

	err := u.repo.CancelDelayRequest(appointID, patientID)
	if err != nil {
		
		if err == entities.ErrNotFound {
			return fmt.Errorf("delay request not found, or it has already been processed")
		}
		return fmt.Errorf("failed to cancel request: %w", err)
	}

	return nil
}