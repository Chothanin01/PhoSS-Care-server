package usecases

import (
	"fmt"

	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"
	"github.com/google/uuid"
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


func (u *appointmentUsecase) CreateAppointment(req *entities.AppointmentCreateReq, adminID uuid.UUID) (*entities.AppointmentRes, error) {
    var result *entities.AppointmentRes

    if req.PatientID == uuid.Nil || req.DiseaseID == uuid.Nil || req.DoctorTitle == "" || req.DoctorFirstName == "" ||
        req.DoctorLastName == "" || req.NextDoctorTitle == "" || req.NextDoctorFirstName == "" || req.NextDoctorLastName == "" ||
        req.Purpose == "" || req.Place == "" || req.Date == "" || req.StartTime == "" || req.EndTime == "" {
        return nil, fmt.Errorf("missing required fields")
    }

    err := u.tx.Do(func(r entities.RepositorySet) error {
        appointRepo := r.AppointmentRepo

        oldAppoint, err := appointRepo.FindOngoing(req.PatientID, req.DiseaseID)
        if err != nil {
            return err
        }

        doctor := req.DoctorTitle + req.DoctorFirstName + " " + req.DoctorLastName

        isVaccine, err := appointRepo.IsVaccineDisease(req.DiseaseID)
        if err != nil {
            return err
        }

        if isVaccine {
            return fmt.Errorf("vaccine appointments must be created using vaccine appointment endpoint")
        }

        if oldAppoint == nil {
            oldEntity := &entities.AppointmentEntity{
                Doctor:    doctor,
                Status:    "completed",
                Symptom:   req.Symptom,
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
            err := appointRepo.UpdateSymptomNote(doctor, oldAppoint.ID, req.Symptom, req.Note, adminID)
            if err != nil {
                return err
            }

            err = appointRepo.CompleteAppoint(oldAppoint.ID, adminID)
            if err != nil {
                return err
            }
        }

        nextDoctor := req.NextDoctorTitle + req.NextDoctorFirstName + " " + req.NextDoctorLastName

		colorStatus := utils.CalculateColorStatus(float64(req.Health.Sugar), float64(req.Health.Pressure))

        newAppoint := &entities.AppointmentEntity{
            Doctor:    nextDoctor,
            Status:    "ongoing",
            Purpose:   req.Purpose,
            Place:     req.Place,
            StartTime: req.StartTime,
            EndTime:   req.EndTime,
            Date:      req.Date,
            PatientID: req.PatientID,
            DiseaseID: req.DiseaseID,
			ColorStatus: colorStatus,
            CreatedBy: adminID,
            UpdatedBy: adminID,
            Health: entities.Health{
                Weight:   req.Health.Weight,
                Height:   req.Health.Height,
                BMI:      req.Health.BMI,
                Pulse:    req.Health.Pulse,
                Pressure: req.Health.Pressure,
                Sugar:    req.Health.Sugar,
            },
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
            Doctor:    savedAppoint.Doctor,
            Purpose:   savedAppoint.Purpose,
            Place:     savedAppoint.Place,
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
	req.VaccineDoctorTitle == "" || req.VaccineDoctorFirstName == "" || req.VaccineDoctorLastName == "" ||
	req.DoctorTitle == "" || req.DoctorFirstName == "" || req.DoctorLastName == "" || req.Place == "" ||
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
			return fmt.Errorf("Patient does not have vaccine disease record.")
		}

		vaccineExists, err := repo.CheckVaccineExists(req.VaccineID)
		if err != nil {
			return err
		}

		if !vaccineExists {
			return fmt.Errorf("vaccine id does not exist")
		}

		if req.OldVaccineID != uuid.Nil {

			oldExists, err := repo.CheckVaccineExists(req.OldVaccineID)
			if err != nil {
				return err
			}

			if !oldExists {
				return fmt.Errorf("old vaccine id does not exist")
			}
		}
		vaccineDiseaseID, err := repo.FindVaccineDiseaseID()
		if err != nil {
			return err
		}

		lastRecord, err := repo.FindOngoingVaccination(req.PatientID)
		vaccineDoctor := req.VaccineDoctorTitle + req.VaccineDoctorFirstName + " " + req.VaccineDoctorLastName
		doctor := req.DoctorTitle + req.DoctorFirstName + " " + req.DoctorLastName
		
		
		if err == nil {

			err = repo.UpdateVaccineDoctor(lastRecord.ID, vaccineDoctor, adminID)
			if err != nil {
				return err
			}

			err = repo.CompleteVaccinationRecord(lastRecord.ID, adminID)
			if err != nil {
				return err
			}

		} else {

			oldAppoint := &entities.AppointmentEntity{
				Purpose: "ฉีดวัคซีน",
				Status:    "completed",
				DiseaseID: vaccineDiseaseID,
				PatientID: req.PatientID,
				Date: 	   req.Date,
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
				req.DoseNumber,
				vaccineDoctor,
				adminID,
			)

			if err != nil {
				return err
			}
		}

		newAppoint := &entities.AppointmentEntity{
			Doctor:    doctor,
			Status:    "ongoing",
			Place:     req.Place,
			Date:      req.Date,
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
		if maxDose == 0 {
			nextDose = 1
		}

		err = repo.CreateVaccinationRecord(
			req.VaccineID,
			saved.ID,
			nextDose,
			doctor,
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
			Doctor:    saved.Doctor,
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
	req.StartTime == "" || req.EndTime == "" || req.DoctorTitle == "" || 
	req.DoctorFirstName == "" || req.DoctorLastName == "" {
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

		doctor := req.DoctorTitle + req.DoctorFirstName + " " + req.DoctorLastName

		updated, err := repo.UpdateAppointment(
			req.AppointID,
			req.Purpose,
			req.Place,
			req.Date,
			req.StartTime,
			req.EndTime,
			doctor,
			adminID,
		)

		if err != nil {
			return err
		}

		result = &entities.AppointmentRes{
			ID: 	   updated.ID,
			No:        updated.No,
			Status:    updated.Status,
			Purpose:   updated.Purpose,
			Place:     updated.Place,
			Date:      updated.Date,
			StartTime: updated.StartTime,
			EndTime:   updated.EndTime,
			Doctor:    updated.Doctor,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *appointmentUsecase) UpdateVaccineAppointment(req *entities.VaccineAppointmentUpdateReq, adminID uuid.UUID) (*entities.VaccineAppointmentRes, error) {

	if req.AppointID == uuid.Nil || req.PatientID == uuid.Nil ||req.VaccineID == uuid.Nil || req.Place == "" || req.Date == "" || 
	req.StartTime == "" || req.EndTime == "" || req.DoctorTitle == "" || req.DoctorFirstName == "" ||req.DoctorLastName == "" {
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
			return fmt.Errorf("Patient does not have vaccine disease record.")
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
		if maxDose == 0 {
			nextDose = 1
		}

		err = repo.UpdateVaccinationRecord(
			req.AppointID,
			req.VaccineID,
			nextDose,
			adminID,
		)
		if err != nil {
			return err
		}

		doctor := req.DoctorTitle + req.DoctorFirstName + " " + req.DoctorLastName

		updated, err := repo.UpdateVaccineAppointment(
			req.AppointID,
			req.Place,
			req.Date,
			req.StartTime,
			req.EndTime,
			doctor,
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
			Doctor:    updated.Doctor,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}