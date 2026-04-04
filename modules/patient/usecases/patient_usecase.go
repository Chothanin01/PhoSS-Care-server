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

func (u *GetPatientUsecase) GetPatientFullInfo(userID uuid.UUID) (*entities.PatientInfoRes, error) {
    patient, err := u.readRepo.GetPatientFullInfo(userID)
    if err != nil {
        return nil, err
    }

    addr := patient.Address
    formattedAddress := fmt.Sprintf("%s หมู่ %s ซอย %s ถนน %s ตำบล %s อำเภอ %s จังหวัด %s %s",
        addr.HouseNumber, addr.VillageNumber, addr.Alley, addr.Road,
        addr.SubDistrict, addr.District, addr.Province, addr.ZipCode)
    
	now := time.Now()
    dob := patient.DOB
    years := now.Year() - dob.Year()
    months := int(now.Month()) - int(dob.Month())
    days := now.Day() - dob.Day()

    if days < 0 {
        months--
    }
    if months < 0 {
        years--
        months += 12
    }

    var weight, height float32
    if len(patient.Appointments) > 0 {
        latest := patient.Appointments[0]
        weight = float32(latest.Health.Weight)
        height = float32(latest.Health.Height)
    } else {
		weight = patient.Weight
        height = patient.Height
    }

    var relWrapper entities.RelativeWrapper
    var offWrapper entities.OfficerWrapper

    for _, r := range patient.Relatives {
        
        fullname := fmt.Sprintf("%s%s %s", r.Title, r.FirstName, r.LastName)
        
        switch r.Role {
        case "kin":
            relWrapper.Kin = &entities.RelativeDetail{
                Fullname:    fullname,
                PhoneNumber: r.PhoneNumber,
                Role:        r.Role,
                Address:     fmt.Sprintf("%s หมู่ %s ซอย %s ถนน %s ตำบล %s อำเภอ %s จังหวัด %s %s",
								addr.HouseNumber, addr.VillageNumber, addr.Alley, addr.Road,
								addr.SubDistrict, addr.District, addr.Province, addr.ZipCode), 
            }
        case "caretaker":
            relWrapper.Caretaker = &entities.RelativeDetail{
                Fullname:    fullname,
                PhoneNumber: r.PhoneNumber,
                Role:        r.Role,
                Address:     fmt.Sprintf("%s หมู่ %s ซอย %s ถนน %s ตำบล %s อำเภอ %s จังหวัด %s %s",
								addr.HouseNumber, addr.VillageNumber, addr.Alley, addr.Road,
								addr.SubDistrict, addr.District, addr.Province, addr.ZipCode),
            }
        case "medicine":
            relWrapper.Medicine = &entities.RelativeDetail{
                Fullname:    fullname,
                PhoneNumber: r.PhoneNumber,
                Role:        r.Role,
                Address:     fmt.Sprintf("%s หมู่ %s ซอย %s ถนน %s ตำบล %s อำเภอ %s จังหวัด %s %s",
								addr.HouseNumber, addr.VillageNumber, addr.Alley, addr.Road,
								addr.SubDistrict, addr.District, addr.Province, addr.ZipCode),
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
			AgeYears:    years,
			AgeMonths:   months,
			AgeDays:     days,
			Allergy:     patient.Allergy,
			PhoneNumber: patient.PhoneNumber,
			Address:     formattedAddress,
			Weight:      weight,
			Height:      height,
			Nationality: patient.Nationality,
			Ethnicity:   patient.Ethnicity,
			DOB:         patient.DOB.Format("2006-01-02"),
		},
		Relative: relWrapper,
        Officer:  offWrapper,
    }, nil
}