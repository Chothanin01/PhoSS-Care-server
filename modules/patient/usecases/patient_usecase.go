package usecases

import (
	"fmt"
	"time"

	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
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

		if patient.DOB.After(now) {
			ageYears, ageMonths, ageDays = 0, 0, 0
		} else {
			years := now.Year() - patient.DOB.Year()
			months := int(now.Month()) - int(patient.DOB.Month())
			days := now.Day() - patient.DOB.Day()

			if days < 0 {
				months--

				daysInPrevMonth := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location()).Day()

				if patient.DOB.Day() > daysInPrevMonth {
					days = now.Day()
				} else {
					days += daysInPrevMonth
				}
			}

			if months < 0 {
				months += 12
				years--
			}

			ageYears = years
			ageMonths = months
			ageDays = days
		}
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

func (u *GetPatientUsecase) GetPatientFullInfo(userID uuid.UUID) (*entities.PatientInfoRes, error) {
    patient, err := u.readRepo.GetPatientFullInfo(userID)
    if err != nil {
        return nil, err
    }

    addr := patient.Address
    formattedAddress := fmt.Sprintf("%s หมู่ %s ซอย %s ถนน %s ตำบล %s อำเภอ %s จังหวัด %s %s",
        addr.HouseNumber, addr.VillageNumber, addr.Alley, addr.Road,
        addr.SubDistrict, addr.District, addr.Province, addr.ZipCode)
    
	var ageYears, ageMonths, ageDays int

	if !patient.DOB.IsZero() {
		now := time.Now()

		if patient.DOB.After(now) {
			ageYears, ageMonths, ageDays = 0, 0, 0
		} else {
			years := now.Year() - patient.DOB.Year()
			months := int(now.Month()) - int(patient.DOB.Month())
			days := now.Day() - patient.DOB.Day()

			if days < 0 {
				months--

				daysInPrevMonth := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location()).Day()

				if patient.DOB.Day() > daysInPrevMonth {
					days = now.Day()
				} else {
					days += daysInPrevMonth
				}
			}

			if months < 0 {
				months += 12
				years--
			}

			ageYears = years
			ageMonths = months
			ageDays = days
		}
	}

    var weight, height, bmi float32
    if len(patient.Appointments) > 0 {
        latest := patient.Appointments[0]
        weight = float32(latest.Health.Weight)
        height = float32(latest.Health.Height)
		bmi = float32(latest.Health.BMI)
    } else {
		weight = patient.Weight
        height = patient.Height
		if height > 0 {
		heightInMeters := height / 100
		
		bmi = weight / (heightInMeters * heightInMeters)
		
	} else {
		bmi = 0 
	}

    }

    var relWrapper entities.RelativeWrapper
    var offWrapper entities.OfficerWrapper

    for _, r := range patient.Relatives {
        
        fullname := fmt.Sprintf("%s%s %s", r.Title, r.FirstName, r.LastName)

		relAddr := r.Address
        
        switch r.Role {
        case "kin":
            relWrapper.Kin = &entities.RelativeDetail{
                Fullname:    fullname,
                PhoneNumber: r.PhoneNumber,
                Role:        r.Role,
                Address:     fmt.Sprintf("%s หมู่ %s ซอย %s ถนน %s ตำบล %s อำเภอ %s จังหวัด %s %s",
								relAddr.HouseNumber, relAddr.VillageNumber, relAddr.Alley, relAddr.Road,
								relAddr.SubDistrict, relAddr.District, relAddr.Province, relAddr.ZipCode), 
            }
        case "caretaker":
            relWrapper.Caretaker = &entities.RelativeDetail{
                Fullname:    fullname,
                PhoneNumber: r.PhoneNumber,
                Role:        r.Role,
                Address:     fmt.Sprintf("%s หมู่ %s ซอย %s ถนน %s ตำบล %s อำเภอ %s จังหวัด %s %s",
								relAddr.HouseNumber, relAddr.VillageNumber, relAddr.Alley, relAddr.Road,
								relAddr.SubDistrict, relAddr.District, relAddr.Province, relAddr.ZipCode),
            }
        case "medicine":
            relWrapper.Medicine = &entities.RelativeDetail{
                Fullname:    fullname,
                PhoneNumber: r.PhoneNumber,
                Role:        r.Role,
                Address:     fmt.Sprintf("%s หมู่ %s ซอย %s ถนน %s ตำบล %s อำเภอ %s จังหวัด %s %s",
								relAddr.HouseNumber, relAddr.VillageNumber, relAddr.Alley, relAddr.Road,
								relAddr.SubDistrict, relAddr.District, relAddr.Province, relAddr.ZipCode),
            }
        case "house":
            offWrapper.House = &entities.OfficerDetail{
                Fullname: fullname,
                Role:     r.Role,
            }
        case "nurse":
            offWrapper.Nurse = &entities.OfficerDetail{
                Fullname: fullname,
                Role:     r.Role,
            }
        }
    }

    return &entities.PatientInfoRes{
        Patient: entities.PatientDetail{
			Fullname:    fmt.Sprintf("%s%s %s", patient.Title, patient.FirstName, patient.LastName),
        	Sex:         patient.Sex,
			IDCard:      patient.IDCard,
			HnNumber:    patient.HnID,
			Rights:      patient.Rights,
			AgeYears:    ageYears,
			AgeMonths:   ageMonths,
			AgeDays:     ageDays,
			Allergy:     patient.Allergy,
			PhoneNumber: patient.PhoneNumber,
			Address:     formattedAddress,
			Weight:      weight,
			Height:      height,
			Nationality: patient.Nationality,
			Ethnicity:   patient.Ethnicity,
			DOB:         patient.DOB.Format("2006-01-02"),
			BMI:         bmi,
		},
		Relative: relWrapper,
        Officer:  offWrapper,
    }, nil
}

func (u *GetPatientUsecase) GetPatientDiseases(patientID uuid.UUID) ([]entities.DiseaseItem, error) {
	
	if patientID == uuid.Nil {
		return nil, fmt.Errorf("patient ID is required")
	}

    diseases, err := u.readRepo.GetPatientDiseases(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch patient diseases: %w", err)
	}

	if diseases == nil {
		diseases = make([]entities.DiseaseItem, 0)
	}

	return diseases, nil
}