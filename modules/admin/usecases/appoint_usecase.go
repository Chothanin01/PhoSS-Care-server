package usecases

import (

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
)

type appointmentUsecase struct {
	tx entities.AppointmentTransaction
	repo entities.AppointmentRepository
}

func NewAppointmentUsecase(tx entities.AppointmentTransaction,repo entities.AppointmentRepository,) entities.AppointmentUsecase {
	return &appointmentUsecase{
		tx:   tx,
		repo: repo,
	}
}


func (u *appointmentUsecase) CreateAppointment(req *entities.AppointmentCreateReq, adminID uuid.UUID) (*entities.AppointmentCreateRes, error) {

	var result *entities.AppointmentCreateRes

	err := u.tx.Do(func(r entities.RepositorySet) error {

		appointRepo := r.AppointmentRepo

		oldAppoint, err := appointRepo.FindOngoing(req.PatientID, req.DiseaseID)
		if err != nil {
			return err
		}

		doctor := req.DoctorTitle + req.DoctorFirstName + " " + req.DoctorLastName

		if oldAppoint == nil {

			oldEntity := &entities.AppointmentEntity{
				Doctor:    doctor,
				Status:    "completed",
				Symtom:    req.Symptom,
				Note:      req.Note,
				PatientID: req.PatientID,
				DiseaseID: req.DiseaseID,
				CreatedBy: adminID,
				UpdatedBy: adminID,
			}

			_, err := appointRepo.CreateAppointment(oldEntity, adminID)
			if err != nil {
				return err
			}
		} else {

			err := appointRepo.UpdateSymptomNote(doctor ,oldAppoint.ID, req.Symptom, req.Note, adminID)
			if err != nil {
				return err
			}

			err = appointRepo.CompleteAppoint(oldAppoint.ID, adminID)
			if err != nil {
				return err
			}
		}


		nextDoctor := req.NextDoctorTitle + req.NextDoctorFirstName + " " + req.NextDoctorLastName

		newAppoint := &entities.AppointmentEntity{
			Doctor:    nextDoctor,
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

			err := appointRepo.CreateHealthRecord(health, req.PatientID, savedAppoint.ID, adminID)
			if err != nil {
				return err
			}
		}

		result = &entities.AppointmentCreateRes{
			ID:      savedAppoint.ID,
			Status:  savedAppoint.Status,
			No:      savedAppoint.No,
			Date:    savedAppoint.Date,
			Time:    savedAppoint.Time,
			Doctor:  savedAppoint.Doctor,
			Purpose: savedAppoint.Purpose,
			Place:   savedAppoint.Place,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}



func (u *appointmentUsecase) FindOngoingVaccination(patientID uuid.UUID) (*entities.VaccineFullDetail, error) {

	record, err := u.repo.FindOngoingVaccination(patientID)
	if err != nil {
		return nil, err
	}

	res := &entities.VaccineFullDetail{
		VaccineID: record.VaccineID,
		Date:      record.CreatedAt.Format("2006-01-02"),
		Type:      record.Vaccine.Type,
		Effect:    record.Vaccine.Effect,
		Note:      record.Vaccine.Note,
		Age:       record.Vaccine.Age,
	}

	return res, nil
}