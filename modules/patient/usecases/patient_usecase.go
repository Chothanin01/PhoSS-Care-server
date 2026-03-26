package usecases

import (
	"fmt"
	"time"

	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"
	"github.com/google/uuid"
)

type GetPatientUsecase struct {
	readRepo entities.GetPatientRepo
}

func NewGetPatientUsecase(readRepo entities.GetPatientRepo) entities.GetPatientUsecase {
	return &GetPatientUsecase{readRepo: readRepo}
}

func (u *GetPatientUsecase) GetPatientBasicInfo(patientID uuid.UUID) (*entities.PatientBasicInfoRes, error) {

	patient, err := u.readRepo.GetPatientBasicInfo(patientID)
	if err != nil {
		return nil, err
	}

	fullname := patient.Title + patient.FirstName + " " + patient.LastName

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

		ageYears = years
		ageMonths = months
		ageDays = days
	}

	res := &entities.PatientBasicInfoRes{
		PatientID: patient.ID,
		FullName:  fullname,
		HnNumber:  patient.HnID,
		AgeYears:  ageYears,
		AgeMonths: ageMonths,
		AgeDays:   ageDays,
	}

	return res, nil
}

func (u *GetPatientUsecase) GetPatientAppointment(patientID uuid.UUID) (*entities.PatientAppointment, error) {
	if patientID == uuid.Nil {
		return nil, fmt.Errorf("user ID cannot be nil")
	}

	appointDB, err := u.readRepo.GetPatientAppointment(patientID)
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
