package usecases

import (
	"time"

	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"
	"github.com/google/uuid"
)

type patientGetUsecase struct {
	readRepo entities.PatientGetRepo
}

func NewPatientGetUsecase(readRepo entities.PatientGetRepo) entities.PatientGetUsecase {
	return &patientGetUsecase{readRepo: readRepo}
}

func (u *patientGetUsecase) GetPatientBasicInfo(userID uuid.UUID) (*entities.PatientBasicInfoRes, error) {

	patient, err := u.readRepo.GetPatientBasicInfo(userID)
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
		AgeYears:     ageYears,
		AgeMonths:    ageMonths,
		AgeDays:      ageDays,
	}

	return res, nil
}