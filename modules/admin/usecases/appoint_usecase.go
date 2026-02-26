package usecases

import (
	"fmt"

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