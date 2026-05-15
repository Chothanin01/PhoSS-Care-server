package usecases

import (
	"fmt"

	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"
	"github.com/google/uuid"
)

type appointmentUsecase struct {
	tx   entities.AppointmentTransaction
	repo entities.AppointmentRepository
}

func NewAppointmentUsecase(tx entities.AppointmentTransaction, repo entities.AppointmentRepository) entities.AppointmentUsecase {
	return &appointmentUsecase{
		tx:   tx,
		repo: repo,
	}
}

func (u *appointmentUsecase) CreateAppointment(req *entities.AppointmentCreateReq, adminID uuid.UUID) (*entities.AppointmentRes, error) {
	var result *entities.AppointmentRes

	if req.PatientID == uuid.Nil || req.DiseaseID == uuid.Nil || req.DoctorID == uuid.Nil || req.NextDoctorID == uuid.Nil ||
		req.Purpose == "" || req.Place == "" || req.Date == "" || req.StartTime == "" || req.EndTime == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	err := u.tx.Do(func(r entities.RepositorySet) error {
		appointRepo := r.AppointmentRepo

		oldAppoint, err := appointRepo.FindOngoing(req.PatientID, req.DiseaseID)
		if err != nil {
			return err
		}

		isVaccine, err := appointRepo.IsVaccineDisease(req.DiseaseID)
		if err != nil {
			return err
		}

		if isVaccine {
			return fmt.Errorf("vaccine appointments must be created using vaccine appointment endpoint")
		}

		disease, err := u.repo.FindDisease(req.DiseaseID)
		if err != nil {
			return err
		}

		colorStatus := "none"

        if disease.Name == "โรคเบาหวาน" || disease.Name == "โรคความดันโลหิตสูง" {
            colorStatus = utils.CalculateColorStatus(float64(req.Health.Sugar), float64(req.Health.Pressure))
        }

		if oldAppoint == nil {

			oldEntity := &entities.AppointmentEntity{
				DoctorID:  req.DoctorID,
				Status:    "completed",
				Symptom:   req.Symptom,
				Note:      req.Note,
				Date:      req.OldDate,
				PatientID: req.PatientID,
				DiseaseID: req.DiseaseID,
				Health:    req.Health,
				ColorStatus: colorStatus,
				CreatedBy: adminID,
				UpdatedBy: adminID,
			}

			_, err := appointRepo.CreateAppointment(oldEntity, adminID)
			if err != nil {
				return err
			}
		} else {
			err := appointRepo.UpdateSymptomNote(req.DoctorID, oldAppoint.ID, colorStatus, req.Symptom, req.Note, adminID)
			if err != nil {
				return err
			}

			err = appointRepo.CompleteAppoint(oldAppoint.ID, adminID)
			if err != nil {
				return err
			}
		}

		newAppoint := &entities.AppointmentEntity{
			DoctorID:    req.NextDoctorID,
			Status:      "ongoing",
			Purpose:     req.Purpose,
			Note: 		 req.Prepare,		
			Place:       req.Place,
			StartTime:   req.StartTime,
			EndTime:     req.EndTime,
			Date:        req.Date,
			PatientID:   req.PatientID,
			DiseaseID:   req.DiseaseID,
			CreatedBy:   adminID,
			UpdatedBy:   adminID,
		}

		savedAppoint, err := appointRepo.CreateAppointment(newAppoint, adminID)
		if err != nil {
			return err
		}

		result = &entities.AppointmentRes{
			ID:        savedAppoint.ID,
			Status:    savedAppoint.Status,
			No:        savedAppoint.No,
			Date:      savedAppoint.Date,
			StartTime: savedAppoint.StartTime,
			EndTime:   savedAppoint.EndTime,
			DoctorID:  savedAppoint.DoctorID,
			Purpose:   savedAppoint.Purpose,
			Place:     savedAppoint.Place,
			ColorStatus: colorStatus,
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
		Name:      record.Vaccine.Name,
		Date:      record.CreatedAt.Format("2006-01-02"),
		Type:      record.Vaccine.Type,
		Effect:    record.Vaccine.Effect,
		Note:      record.Vaccine.Note,
		Age:       record.Vaccine.Age,
	}

	return res, nil
}

func (u *appointmentUsecase) CreateVaccineAppointment(req *entities.VaccineAppointmentCreateReq, adminID uuid.UUID) (*entities.VaccineAppointmentRes, error) {
	var result *entities.VaccineAppointmentRes

	if req.PatientID == uuid.Nil || req.VaccineID == uuid.Nil || req.DoseNumber == 0 ||
		req.VaccineDoctorID == uuid.Nil || req.DoctorID == uuid.Nil || req.Place == "" ||
		req.NextDate == "" || req.StartTime == "" || req.EndTime == "" {
		return nil, fmt.Errorf("missing required fields for vaccine appointment")
	}

	err := u.tx.Do(func(r entities.RepositorySet) error {
		repo := r.AppointmentRepo

		hasDisease, err := repo.FindPatientVaccine(req.PatientID)
		if err != nil {
			return err
		}

		if !hasDisease {
			return fmt.Errorf("patient does not have vaccine disease record")
		}

		vaccineExists, err := repo.CheckVaccineExists(req.VaccineID)
		if err != nil {
			return err
		}

		if !vaccineExists {
			return fmt.Errorf("vaccine id does not exist")
		}

		vaccineDiseaseID, err := repo.FindVaccineDiseaseID()
		if err != nil {
			return err
		}

		lastRecord, err := repo.FindOngoingVaccination(req.PatientID)

		if err == nil && lastRecord != nil {
			err = repo.CompleteAppoint(lastRecord.AppointID, adminID)
			if err != nil {
				return err
			}
			
			err = repo.UpdateVaccineDoctor(lastRecord.ID, req.VaccineDoctorID, adminID)
			if err != nil {
				return err
			}

			err = repo.CompleteVaccinationRecord(lastRecord.ID, adminID)
			if err != nil {
				return err
			}
		} else {
			oldAppoint := &entities.AppointmentEntity{
				Purpose:   "ฉีดวัคซีน",
				Status:    "completed",
				DiseaseID: vaccineDiseaseID,
				PatientID: req.PatientID,
				DoctorID:  req.VaccineDoctorID,
				Date:      req.Date,
				CreatedBy: adminID,
				UpdatedBy: adminID,
			}

			old, err := repo.CreateAppointment(oldAppoint, adminID)
			if err != nil {
				return err
			}

			err = repo.CreateVaccinationRecord(
				req.OldVaccineID,
				old.ID,
				"completed",
				req.DoseNumber,
				req.VaccineDoctorID,
				adminID,
			)
			if err != nil {
				return err
			}
		}

		newAppoint := &entities.AppointmentEntity{
			DoctorID:  req.DoctorID,
			Status:    "ongoing",
			Place:     req.Place,
			Date:      req.NextDate,
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			Purpose:   "ฉีดวัคซีน",
			DiseaseID: vaccineDiseaseID,
			PatientID: req.PatientID,
			CreatedBy: adminID,
			UpdatedBy: adminID,
		}

		saved, err := repo.CreateAppointment(newAppoint, adminID)
		if err != nil {
			return err
		}

		maxDose, err := repo.FindMaxDose(req.PatientID, req.VaccineID)
		if err != nil {
			return err
		}

		nextDose := maxDose + 1

		err = repo.CreateVaccinationRecord(
			req.VaccineID,
			saved.ID,
			"ongoing",
			nextDose,
			req.DoctorID,
			adminID,
		)
		if err != nil {
			return err
		}

		result = &entities.VaccineAppointmentRes{
			AppointID: saved.ID,
			No:        saved.No,
			Status:    saved.Status,
			Date:      saved.Date.String(),
			StartTime: saved.StartTime,
			EndTime:   saved.EndTime,
			DoctorID:  saved.DoctorID,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *appointmentUsecase) UpdateAppointment(req *entities.AppointmentUpdateReq, adminID uuid.UUID) (*entities.AppointmentRes, error) {
	if req.AppointID == uuid.Nil || req.Purpose == "" || req.Place == "" || req.Date == "" ||
		req.StartTime == "" || req.EndTime == "" || req.DoctorID == uuid.Nil {
		return nil, fmt.Errorf("missing required fields")
	}

	var result *entities.AppointmentRes

	err := u.tx.Do(func(r entities.RepositorySet) error {
		repo := r.AppointmentRepo

		appoint, err := repo.FindByID(req.AppointID)
		if err != nil {
			return fmt.Errorf("appointment not found")
		}

		if appoint.Status == "completed" {
			return fmt.Errorf("completed appointment cannot be updated")
		}

		isVaccine, err := repo.IsVaccineDisease(appoint.DiseaseID)
		if err != nil {
			return err
		}

		if isVaccine {
			return fmt.Errorf("vaccine appointment cannot be updated in this module")
		}

		updated, err := repo.UpdateAppointment(
			req.AppointID,
			req.Purpose,
			req.Place,
			req.Date,
			req.StartTime,
			req.EndTime,
			req.DoctorID,
			adminID,
		)

		if err != nil {
			return err
		}

		result = &entities.AppointmentRes{
			ID:        updated.ID,
			No:        updated.No,
			Status:    updated.Status,
			Purpose:   updated.Purpose,
			Place:     updated.Place,
			Date:      updated.Date,
			StartTime: updated.StartTime,
			EndTime:   updated.EndTime,
			DoctorID:  updated.DoctorID,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *appointmentUsecase) UpdateVaccineAppointment(req *entities.VaccineAppointmentUpdateReq, adminID uuid.UUID) (*entities.VaccineAppointmentRes, error) {
	if req.AppointID == uuid.Nil || req.PatientID == uuid.Nil || req.VaccineID == uuid.Nil || req.Place == "" || req.Date == "" ||
		req.StartTime == "" || req.EndTime == "" || req.DoctorID == uuid.Nil {
		return nil, fmt.Errorf("missing required fields")
	}

	var result *entities.VaccineAppointmentRes

	err := u.tx.Do(func(r entities.RepositorySet) error {
		repo := r.AppointmentRepo

		appoint, err := repo.FindByID(req.AppointID)
		if err != nil {
			return fmt.Errorf("appointment not found")
		}

		if appoint.Status == "completed" {
			return fmt.Errorf("completed appointment cannot be updated")
		}

		if len(appoint.Vaccinations) == 0 {
			return fmt.Errorf("this appoint is not vaccine appoint")
		}

		hasDisease, err := repo.FindPatientVaccine(req.PatientID)
		if err != nil {
			return err
		}

		if !hasDisease {
			return fmt.Errorf("patient does not have vaccine disease record")
		}

		exists, err := repo.CheckVaccineExists(req.VaccineID)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("vaccine not found")
		}

		maxDose, err := repo.FindMaxDose(req.PatientID, req.VaccineID)
		if err != nil {
			return err
		}

		nextDose := maxDose + 1

		err = repo.UpdateVaccinationRecord(
			req.AppointID,
			req.VaccineID,
			nextDose,
			adminID,
		)
		if err != nil {
			return err
		}

		updated, err := repo.UpdateVaccineAppointment(
			req.AppointID,
			req.Place,
			req.Date,
			req.StartTime,
			req.EndTime,
			req.DoctorID,
			adminID,
		)
		if err != nil {
			return err
		}

		result = &entities.VaccineAppointmentRes{
			AppointID: updated.ID,
			No:        updated.No,
			Status:    updated.Status,
			Date:      updated.Date.Format("2006-01-02"),
			StartTime: updated.StartTime,
			EndTime:   updated.EndTime,
			DoctorID:  updated.DoctorID,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *appointmentUsecase) FindAllDoctor(role string) ([]entities.Doctor, error) {
    
    doctors, err := u.repo.FindAllDoctor(role)
    if err != nil {
        return nil, err
    }

    return doctors, nil
}