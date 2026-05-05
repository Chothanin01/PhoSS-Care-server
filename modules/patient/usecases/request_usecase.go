package usecases

import (
	"fmt"
	"time"
	"strings"

	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
	"github.com/google/uuid"
)

type requestCommandUsecase struct {
	repo entities.RequestCommandRepo
}

func NewRequestCommandUsecase(repo entities.RequestCommandRepo) entities.RequestCommandUsecase {
	return &requestCommandUsecase{repo: repo}
}

type requestQueryUsecase struct {
	repo entities.RequestQueryRepo
}

func NewRequestQueryUsecase(repo entities.RequestQueryRepo) entities.RequestQueryUsecase {
	return &requestQueryUsecase{repo: repo}
}

func (u *requestCommandUsecase) SubmitDocumentRequest(userID uuid.UUID, patientID uuid.UUID, payload *entities.CreateDocumentReq) error {
	if len(payload.Requests) == 0 {
		return fmt.Errorf("request payload cannot be empty")
	}

	var requestsToSave []entities.RequestEntity
	var notificationsToSave []entities.NotificationEntity
	now := time.Now()

	for _, item := range payload.Requests {
		
		if item.Type == "medical" {
			if item.DiseaseID == uuid.Nil {
				return fmt.Errorf("disease_id is required for medical requests")
			}

			appointID, err := u.repo.GetLatestAppointID(patientID, item.DiseaseID)
			if err != nil {
				return fmt.Errorf("failed to process medical request: %w", err)
			}

			diseaseName, err := u.repo.GetDiseaseName(item.DiseaseID)
			if err != nil {
				return fmt.Errorf("failed to find disease info: %w", err)
			}

			diseaseIDCopy := item.DiseaseID 
			requestsToSave = append(requestsToSave, entities.RequestEntity{
				RequestType: "medical",
				Description: diseaseName,
				Date:        now,
				Status:      "pending",
				PatientID:   patientID,
				AppointID:   appointID,
				DiseaseID:   &diseaseIDCopy,
				CreatedBy:   userID,
			})
			notificationsToSave = append(notificationsToSave, entities.NotificationEntity{
				Header:    entities.NotiHeaderDocument, 
				Body:      entities.NotiBodyDocumentSent,
				PatientID: patientID,
				CreatedBy: userID,
			})
		}

		if item.Type == "document" {
			if len(item.DocumentTypes) == 0 {
				return fmt.Errorf("document_types array cannot be empty for document requests")
			}

			joinedDocs := strings.Join(item.DocumentTypes, ", ")

			requestsToSave = append(requestsToSave, entities.RequestEntity{
				RequestType: "document",
				Description: joinedDocs,
				Date:        now,
				Status:      "pending",
				PatientID:   patientID,
				AppointID:   nil,
				DiseaseID:   nil,
				CreatedBy:   userID,
			})

			notificationsToSave = append(notificationsToSave, entities.NotificationEntity{
				Header:    entities.NotiHeaderMedical, 
				Body:      entities.NotiBodyMedicalSent,
				PatientID: patientID,
				CreatedBy: userID,
			})
		}
	}

	return u.repo.SaveMultipleRequests(requestsToSave, notificationsToSave)
}

func (u *requestQueryUsecase) GetAvailableDocumentOptions(patientID uuid.UUID) ([]entities.AvailableRequestOption, error) {

	hasAppoint, _ := u.repo.CheckHasAnyAppointment(patientID)
	hasVaccine, _ := u.repo.CheckHasVaccineHistory(patientID)
	_, latestDiseaseID, _ := u.repo.GetLatestCompletedAppoint(patientID)

	pendingReqs, _ := u.repo.GetPendingRequests(patientID)

	isDocPending := func(docName string) bool {
		for _, req := range pendingReqs {
			if req.RequestType == "document" {

				requestedDocs := strings.Split(req.Description, ",")
				for _, requestedDoc := range requestedDocs {
						
					if strings.TrimSpace(requestedDoc) == docName {
						return true
					}
				}
			}
		}
		return false
	}
	

	isMedicalPending := func(targetDiseaseID *uuid.UUID) bool {
		if targetDiseaseID == nil {
			return false
		}
		for _, req := range pendingReqs {
			if req.RequestType == "medical" && req.DiseaseID != nil {
				if *req.DiseaseID == *targetDiseaseID {
					return true
				}
			}
		}
		return false
	}

	var options []entities.AvailableRequestOption

	options = append(options, entities.AvailableRequestOption{
		Name:      "ข้อมูลผู้ป่วย",
		Type:      "document",
		Available: true,
		Disabled:  isDocPending("ข้อมูลผู้ป่วย"),
	})

	options = append(options, entities.AvailableRequestOption{
		Name:      "ประวัติการรักษา",
		Type:      "document",
		Available: hasAppoint,
		Disabled:  isDocPending("ประวัติการรักษา"),
	})

	options = append(options, entities.AvailableRequestOption{
		Name:      "ประวัติการฉีดวัคซีน",
		Type:      "document",
		Available: hasVaccine,
		Disabled:  isDocPending("ประวัติการฉีดวัคซีน"),
	})

	isCertAvailable := latestDiseaseID != nil
	options = append(options, entities.AvailableRequestOption{
		Name:      "ใบรับรองแพทย์",
		Type:      "medical",
		Available: isCertAvailable,
		Disabled:  isMedicalPending(latestDiseaseID),
		DiseaseID: latestDiseaseID,
	})

	return options, nil
}

func (u *vaccineQueryUsecase) GetVaccineDetail(patientID uuid.UUID, vaccineID uuid.UUID) (*entities.VaccineDetailEntity, error) {
	
	if patientID == uuid.Nil {
		return nil, fmt.Errorf("patient ID is required")
	}
	if vaccineID == uuid.Nil {
		return nil, fmt.Errorf("vaccine ID is required")
	}

	detail, err := u.repo.GetVaccineDetail(patientID, vaccineID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vaccine details: %w", err)
	}

	return detail, nil
}