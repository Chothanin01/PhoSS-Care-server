package usecases

import (
	"fmt"
    "time"

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
)

type appointmentUsecase struct {
	tx entities.AppointmentTransaction
}

func NewAppointmentUsecase(tx entities.AppointmentTransaction) entities.AppointmentUsecase {
	return &appointmentUsecase{tx: tx}
}


func (u *appointmentUsecase) CreateAppointment(req *entities.AppointmentCreateReq, adminID uuid.UUID) (*entities.AppointmentCreateRes, error) {
	var result *entities.AppointmentCreateRes

	err := u.tx.Do(func(r entities.RepositorySet) error {
		appointRepo := r.AppointmentRepo

		if req.DoctorFirstName == "" || req.DoctorLastName == "" || req.Place == "" ||
			req.Note == "" || req.Health.Weight == 0 || req.Health.Height == 0 || req.Symptom == "" ||
			req.PatientID == uuid.Nil || req.DiseaseID == uuid.Nil || req.Purpose == "" ||
			req.Date == "" || req.Time == "" {
			return fmt.Errorf("missing required fields for appointment creation")
		}

        exists, err := appointRepo.DiseaseExists(req.DiseaseID)
		if err != nil {
			return fmt.Errorf("failed to verify disease: %w", err)
		}
		if !exists {
			return fmt.Errorf("invalid disease_id: disease not found")
		}

		oldAppoint, err := appointRepo.FindOngoing(req.PatientID, req.DiseaseID)
		if err != nil {
			return fmt.Errorf("find ongoing appointment: %w", err)
		}

		if oldAppoint != nil {
			if err := appointRepo.CompleteAppoint(oldAppoint.ID, adminID); err != nil {
				return fmt.Errorf("complete ongoing appointment: %w", err)
			}
		}

		newAppoint := &entities.AppointmentEntity{
			Doctor:    req.DoctorTitle + req.DoctorFirstName + " " + req.DoctorLastName,
			Status:    "ongoing",
			Purpose:   req.Purpose,
			Place:     req.Place,
			Time:      req.Time,
			Date:      req.Date,
			PatientID: req.PatientID,
			DiseaseID: req.DiseaseID,
			CreatedBy: adminID,
			UpdatedBy: adminID,
		}

		savedAppoint, err := appointRepo.CreateAppointment(newAppoint, adminID)
		if err != nil {
			return fmt.Errorf("create appointment: %w", err)
		}

		if req.Health.Height > 0 && req.Health.Weight > 0 {
			health := &entities.Health{
				Weight: req.Health.Weight,
				Height: req.Health.Height,
				BMI:    req.Health.BMI,
				Pulse:  req.Health.Pulse,
				Sugar:  req.Health.Sugar,
			}
			if err := appointRepo.CreateHealthRecord(health, req.PatientID, savedAppoint.ID, adminID); err != nil {
				return fmt.Errorf("create health record: %w", err)
			}
		}

		result = &entities.AppointmentCreateRes{
			ID:     	savedAppoint.ID,
			Status: 	savedAppoint.Status,
            No:     	savedAppoint.No,
			Date:   	savedAppoint.Date,
			Time:   	savedAppoint.Time,
			Doctor: 	savedAppoint.Doctor,
			Purpose:  	savedAppoint.Purpose,
			Place:  	savedAppoint.Place,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *appointmentUsecase) UpdateAppointment(req *entities.AppointmentUpdateReq, adminID uuid.UUID) (*entities.AppointmentUpdateRes, error) {
	var result *entities.AppointmentUpdateRes

	err := u.tx.Do(func(r entities.RepositorySet) error {
		appointRepo := r.AppointmentRepo

		appoint, err := appointRepo.FindByID(req.AppointID)

		if appoint.Status != "ongoing" && appoint.Status != "delay" {
			return fmt.Errorf("appointment cannot be updated because its status is '%s'", appoint.Status)
		}
		
		if err != nil {
			return fmt.Errorf("appointment not found: %w", err)
		}

		if req.DoctorFirstName == "" || req.Place == "" || req.Purpose == "" || req.Date == "" || req.Time == "" {
			return fmt.Errorf("missing required fields")
		}

		dateParsed, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return fmt.Errorf("invalid date format (expected YYYY-MM-DD)")
		}

		appoint.Doctor = req.DoctorTitle + req.DoctorFirstName + " " + req.DoctorLastName
		appoint.Symptom = req.Symptom
		appoint.Purpose = req.Purpose
		appoint.Place = req.Place
		appoint.Time = req.Time
		appoint.Date = dateParsed
		appoint.UpdatedBy = &adminID

		if err := appointRepo.UpdateAppointment(appoint); err != nil {
			return err
		}

		health := &entities.Health{
			Weight: req.Health.Weight,
			Height: req.Health.Height,
			BMI:    req.Health.BMI,
			Pulse:  req.Health.Pulse,
			Sugar:  req.Health.Sugar,
		}

		if err := appointRepo.UpdateHealth(req.AppointID, health, adminID); err != nil {
			return err
		}

		result = &entities.AppointmentUpdateRes{
			ID: appoint.ID,
			Doctor: appoint.Doctor,
			Status: appoint.Status,
			Purpose: appoint.Purpose,
			Place: appoint.Place,
			Date: appoint.Date.Format("2006-01-02"),
			Time: appoint.Time,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

