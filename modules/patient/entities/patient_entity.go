package entities

import (
	"errors"
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases" 
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized action")
	ErrInternalError = errors.New("internal system error")
)

type PatientBasicInfoRes struct {
	PatientID uuid.UUID `json:"patient_id"`
	FullName  string    `json:"fullname"`
	HnNumber  string    `json:"hn_number"`
	AgeYears  int       `json:"age_years"`
	AgeMonths int       `json:"age_months"`
	AgeDays   int       `json:"age_days"`
}

type PatientAppointment struct {
	PatientID 		uuid.UUID 				`json:"patient_id"`
	FullName  		string    				`json:"fullname"`
	HnNumber  		string    				`json:"hn_number"`
	AgeYears  		int       				`json:"age_years"`
	AgeMonths 		int       				`json:"age_months"`
	AgeDays   		int       				`json:"age_days"`
	Appointments 	[]AppointmentEntity 	`json:"appoint"`
}

type PatientInfoRes struct {
    Patient  PatientDetail    `json:"patient"`
    Relative RelativeWrapper  `json:"relative"`
    Officer  OfficerWrapper   `json:"officer"`
}

type PatientDetail struct {
	Fullname    		string             `json:"fullname"`
	Sex		 			string             `json:"sex"`
	IDCard      		string             `json:"idcard"`
	HnNumber    		string             `json:"hnnumber"`
	Rights	    		string             `json:"rights"`
	AgeYears    		int                `json:"age_years"`
	AgeMonths   		int                `json:"age_months"`
	AgeDays     		int                `json:"age_days"`
	Allergy     		string             `json:"allergy"`
	PhoneNumber 		string             `json:"phone_number"`
	Address     		string             `json:"address"`
	Weight      		float32            `json:"weight"`
	Height      		float32            `json:"height"`
	Nationality  		string     		   `json:"nationality"`
	Ethnicity    		string       	   `json:"ethnicity"`
	DOB          		string  		   `json:"dob"`    
}

type AddressDetails struct {
	HouseNumber   string `json:"house_number"`
	VillageNumber string `json:"village_number"`
	Alley         string `json:"alley"`
	Road          string `json:"road"`
	SubDistrict   string `json:"subdistrict"`
	District      string `json:"district"`
	Province      string `json:"province"`
	ZipCode       string `json:"zipcode"`
}

type RelativeWrapper struct {
    Kin       *RelativeDetail `json:"kin"`
    Caretaker *RelativeDetail `json:"caretaker"`
    Medicine  *RelativeDetail `json:"medicine"`
}

type OfficerWrapper struct {
    House *OfficerDetail `json:"house"`
    Nurse *OfficerDetail `json:"nurse"`
}

type RelativeDetail struct {
    Fullname    string      `json:"fullname"`
    PhoneNumber string      `json:"phonenumber"`
    Role        string      `json:"role"`
    Address     string 		`json:"address"`
}

type OfficerDetail struct {
    Fullname string `json:"fullname"`
    Role     string `json:"role"`
}

type DiseaseItem struct {
	DiseaseID uuid.UUID `json:"disease_id"`
	Name      string    `json:"name"`
}

type GetPatientRepo interface {
	GetPatientBasicInfo(patientID uuid.UUID) (*databases.Patient, error)
	GetPatientFullInfo(userID uuid.UUID) (*databases.Patient, error)

	GetPatientDiseases(patientID uuid.UUID) ([]DiseaseItem, error)
}

type GetPatientUsecase interface {
	GetPatientBasicInfo(patientID uuid.UUID) (*PatientBasicInfoRes, error)
	GetPatientFullInfo(userID uuid.UUID) (*PatientInfoRes, error)

	GetPatientDiseases(patientID uuid.UUID) ([]DiseaseItem, error)
}