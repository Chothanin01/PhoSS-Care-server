package usecases

import (
	"fmt"
	
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
)

type requestGetUsecase struct {
	repo entities.RequestGetRepo
}

func NewRequestGetUsecase(repo entities.RequestGetRepo) entities.RequestGetUsecase {
	return &requestGetUsecase{repo: repo}
}

type requestUpdateUsecase struct {
	repo entities.RequestUpdateRepo
}

func NewRequestUpdateUsecase(repo entities.RequestUpdateRepo) entities.RequestUpdateUsecase {
	return &requestUpdateUsecase{repo: repo}
}

func (u *requestGetUsecase) GetRequestsWithFilter(params entities.RequestQueryParams) (*entities.RequestListRes, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	data, err := u.repo.GetRequestsWithFilter(params)
	if err != nil {
		return nil, err
	}

	total, err := u.repo.CountRequestsWithFilter(params)
	if err != nil {
		return nil, err
	}

	return &entities.RequestListRes{
		Page:       params.Page,
		PerPage:    params.Limit,
		TotalPages: int((total + int64(params.Limit) - 1) / int64(params.Limit)),
		Data:       data,
	}, nil
}

func (u *requestGetUsecase) GetRequestInfoByID(id uuid.UUID) (*entities.RequestInfoRes, error) {

	req, err := u.repo.GetRequestInfoByID(id)
	if err != nil {
		return nil, err
	}

	p := req.Patient
	fullname := p.Title + p.FirstName + " " + p.LastName

	res := &entities.RequestInfoRes{
		RequestID:   req.ID,
		RequestType: req.RequestType,
		Status:      req.Status,
		PatientID:   p.ID,
		FullName:    fullname,
		IDCard:      p.IDCard,
		HnNumber:    p.HnID,
	}

	res.DiseaseName = req.Disease.Name

	switch req.RequestType {

	case "appoint":
		res.Date = req.Date.Format("2006-01-02")
		res.StartTime = req.StartTime
		res.EndTime = req.EndTime
		res.Doctor = req.Appoint.Doctor.Title + req.Appoint.Doctor.FirstName + " " + req.Appoint.Doctor.LastName 
		res.AppointDate = req.Appoint.Date.Format("2006-01-02")
		res.AppointStartTime = req.Appoint.StartTime
		res.AppointEndTime = req.Appoint.EndTime
		res.AppointID = req.Appoint.ID
		res.Description = req.Description

	case "medical":
		res.CreatedDate = req.CreatedAt.Format("2006-01-02")
		res.Description = req.Description
		res.Doctor = req.Appoint.Doctor.Title + req.Appoint.Doctor.FirstName + " " + req.Appoint.Doctor.LastName 

	case "document":
		res.CreatedDate = req.CreatedAt.Format("2006-01-02")
		res.Description = req.Description
	}

	return res, nil
}

func (u *requestUpdateUsecase) UpdateRequestStatus(req *entities.RequestStatusUpdateReq, adminID uuid.UUID) (*entities.RequestStatusUpdateRes, error) {
	
	request, err := u.repo.FindRequestByID(req.RequestID)

	if err != nil {
		return nil, fmt.Errorf("request not found: %w", err)
	}

	if request.Status != "pending" {
		return nil, fmt.Errorf("cannot update request with status '%s'", request.Status)
	}

	var noti entities.NotificationEntity
	noti.PatientID = request.PatientID
	noti.CreatedBy = adminID
	noti.RequestID = &req.RequestID 

	switch request.RequestType {
	case "appoint":
		noti.Header = "การเลื่อนนัด"
		noti.AppointID = &request.AppointID 

		switch req.Status {
		case "accepted":
			if request.AppointID == uuid.Nil {
				return nil, fmt.Errorf("appointment not linked to this request")
			}
			
			noti.Body = "ระบบได้ยืนยันการเลื่อนนัดของคุณแล้ว"
			err := u.repo.UpdateAppointForAccepted(
				req.RequestID,
				request.AppointID,
				request.Date,
				request.StartTime,
				request.EndTime,
				adminID,
				noti,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to update request: %w", err)
			}

		case "canceled":
			if req.Description == "" {
				return nil, fmt.Errorf("description is required when declining")
			}
			
			noti.Body = "ระบบได้ปฏิเสธการเลื่อนนัดของคุณกรุณาเลื่อนนัดใหม่อีกครั้ง"

			if err := u.repo.UpdateRequestStatus(req.RequestID, "canceled", req.Description, adminID, noti); err != nil {
				return nil, fmt.Errorf("failed to update request: %w", err)
			}

		default:
			return nil, fmt.Errorf("invalid status: must be 'accepted' or 'declined'")
		}

	case "medical", "document":
		if req.Status != "accepted" {
			return nil, fmt.Errorf("only 'accepted' status allowed for %s requests", request.RequestType)
		}

		if request.RequestType == "medical" {
			noti.Header = "ใบรับรองแพทย์"
		} else {
			noti.Header = "เอกสารรับรอง"
		}
		
		noti.Body = "โรงพยาบาลได้เตรียมเอกสารของคุณเเล้ว"

		if err := u.repo.UpdateRequestStatus(req.RequestID, "accepted", request.Description, adminID, noti); err != nil {
			return nil, fmt.Errorf("failed to update request: %w", err)
		}

	default:
		return nil, fmt.Errorf("unknown request type: %s", request.RequestType)
	}

	return &entities.RequestStatusUpdateRes{
		ID:          req.RequestID,
		RequestType: request.RequestType,
		Status:      req.Status,
	}, nil
}