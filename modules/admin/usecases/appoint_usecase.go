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
            req.Note == "" || req.Health.Weight == 0 || req.PatientID == uuid.Nil || req.DiseaseID == uuid.Nil || 
            req.Date == "" || req.Time == "" || req.Health.Height == 0 {
            return fmt.Errorf("missing required fields for follow-up appointment")
        }

		newAppoint := &entities.AppointmentEntity{
			Doctor:    req.DoctorTitle + req.DoctorFirstName + " " + req.DoctorLastName,
			Status:    "ongoing",
			Note:      req.Note,
			Place:     req.Place,
            Time:      req.Time,
            Date:      req.Date,
			PatientID: req.PatientID,
			DiseaseID: req.DiseaseID,
			CreatedBy: adminID,
			UpdatedBy: adminID,
		}

        oldAppoint, err := appointRepo.FindOngoing(req.PatientID, req.DiseaseID)
        if err != nil {
            return fmt.Errorf("find ongoing appointment: %w", err)
        }

        if err := appointRepo.CompleteAppoint(oldAppoint.ID, adminID); err != nil {
            return fmt.Errorf("complete ongoing appointment: %w", err)
        }

		savedAppoint, err := appointRepo.CreateAppointment(newAppoint, adminID)
		if err != nil {
			return err
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
				return err
			}
		}

		result = &entities.AppointmentCreateRes{
			ID:     savedAppoint.ID,
			Status: savedAppoint.Status,
			Date:   savedAppoint.Date,
			Time:   savedAppoint.Time,
			Doctor: savedAppoint.Doctor,
			Note:   savedAppoint.Note,
			Place:  savedAppoint.Place,
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
		if err != nil {
			return fmt.Errorf("appointment not found: %w", err)
		}

		if req.DoctorFirstName == "" || req.Place == "" || req.Note == "" || req.Date == "" || req.Time == "" {
			return fmt.Errorf("missing required fields")
		}

		dateParsed, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return fmt.Errorf("invalid date format (expected YYYY-MM-DD)")
		}

		appoint.Doctor = req.DoctorTitle + req.DoctorFirstName + " " + req.DoctorLastName
		appoint.Symptom = req.Symptom
		appoint.Note = req.Note
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
			Note: appoint.Note,
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