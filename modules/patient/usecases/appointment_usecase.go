package usecases

import (

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
)

type GetAppointmentUsecase struct {
	readRepo entities.GetAppointmentRepo
}

func NewAppointmentUsecase(readRepo entities.GetAppointmentRepo) entities.GetAppointmentUsecase {
	return &GetAppointmentUsecase{readRepo: readRepo}
}

func (u *GetAppointmentUsecase) GetAppointmentDetail(patientID uuid.UUID, diseaseID uuid.UUID) (*entities.AppointmentEntity, error) {
	appointDB, adminDB, err := u.readRepo.GetAppointmentDetail(patientID, diseaseID)
	if err != nil {
		return nil, err
	}

	res := &entities.AppointmentEntity{
		ID:        appointDB.ID,
		No: 	   appointDB.No,	
		Doctor:    appointDB.Doctor,
		Status:    appointDB.Status,
		Purpose:   appointDB.Purpose,
		Place:     appointDB.Place,
		Date:      appointDB.Date.Format("2006-01-02"),
		StartTime: appointDB.StartTime,
		EndTime:   appointDB.EndTime,
		Symptom:   appointDB.Symptom,
		Note:      appointDB.Note,
		Delay:     appointDB.Delay, 
		DiseaseID: appointDB.DiseaseID,
		CreatedAt: appointDB.CreatedAt.Format("2006-01-02"),
	}

	if adminDB != nil && adminDB.FirstName != "" {
        res.CreatedBy = adminDB.Title + adminDB.FirstName + " " + adminDB.LastName
    }

	return res, nil
}