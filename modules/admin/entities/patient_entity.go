package entities

import (

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases" 
)
type RepositorySet struct {
	UserRepo     UserRepository
	PatientRepo  PatientCreateRepo
	RelativeRepo RelativeCreateRepo
	DiseaseRepo  DiseaseRepo
	PateintUpdateRepo PatientUpdateRepo
	AppointmentRepo AppointmentRepository
}

type PatientFullCreateReq struct {
	User struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	} `json:"user"`

	Patient struct {
		Title        string     `json:"title"`
		FirstName    string     `json:"firstname"`
		LastName     string     `json:"lastname"`
		Sex		  	 string	    `json:"sex"`
		DOB          string     `json:"dob"`
		IDCard       string     `json:"idcard"`
		Rights       string     `json:"rights"`
		Nationality  string     `json:"nationality"`
		Ethnicity    string     `json:"ethnicity"`
		PhoneNumber  string     `json:"phonenumber"`
		Weight	     float32    `json:"weight"`
		Height	     float32    `json:"height"`
		Address      Address `json:"address"`
		Allergy      string     `json:"allergy"`
		Diseases     []Disease  `json:"diseases"`
	} `json:"patient"`

	Relative struct {
		Kin       RelativeCreate `json:"kin"`
		Caretaker RelativeCreate `json:"caretaker"`
		Medicine  RelativeCreate `json:"medicine"`
	} `json:"relative"`

	Officer struct {
		House RelativeCreate `json:"house"`
		Nurse RelativeCreate `json:"nurse"`
	} `json:"officer"`

	CreatedBy  uuid.UUID  `json:"created_by"`
}
type PatientCreateReq struct {
	Title        string     `json:"title"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Sex		  	 string     `json:"sex"`
	Dob          string     `json:"dob"`
	HnID         string     `json:"hn_id"`
	IDCard       string     `json:"id_card"`
	Rights       string     `json:"rights"`
	Nationality  string     `json:"nationality"`
	Ethnicity    string     `json:"ethnicity"`
	PhoneNumber  string     `json:"phone_number"`
	Weight	     float32    `json:"weight"`
	Height	     float32    `json:"height"`
	Address      Address `json:"address"`
	Allergy      string     `json:"allergy"`
	Diseases     []Disease `json:"diseases"`
	UserID       uuid.UUID  `json:"user_id"`
}

type PatientCreateRes struct {
	Id          uuid.UUID `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Sex		 	string `json:"sex"`
	HnID        string `json:"hn_id"`
	IDCard      string `json:"id_card"`
	PhoneNumber string `json:"phone_number"`
	Rights      string `json:"rights"`
	Nationality string `json:"nationality"`
	Ethnicity   string `json:"ethnicity"`
}
type Address struct {
	HouseNumber string `json:"house_number"`
	VillageNumber         string `json:"village_number"`
	Alley       string `json:"alley"`
	Road        string `json:"road"`
	SubDistrict string `json:"subdistrict"`
	District    string `json:"district"`
	Province    string `json:"province"`
	ZipCode     string `json:"zipcode"`
}

type PatientListRes struct {
	Success     bool           `json:"success"`
	Message     string         `json:"message"`
	Page        int            `json:"page"`
	PerPage     int            `json:"per_page"`
	TotalPages  int            `json:"total_pages"`
	Data        []PatientHomeInfo  `json:"data"`
}
type PatientHomeInfo struct {
	ID        uuid.UUID          `json:"id"`
	FullName  string             `json:"fullname"`
	IDCard    string             `json:"idcard"`
	HnNumber  string             `json:"hnnumber"`
	Diseases  []DiseaseWithStatus `json:"diseases"`
}

type PatientQueryParams struct {
	Search   string   `query:"search"`
	Diseases []string `query:"disease"`
	Appoint  *bool    `query:"appoint"`
	Overdue  *bool    `query:"overdue"`
	Page     int      `query:"page"`
	Limit    int      `query:"limit"`
}

type PatientInfoRes struct {
	Success     bool           		`json:"success"`
	Message     string         		`json:"message"`
	Data        []PatientData   `json:"data"`
}

type PatientData struct {
	Patient  PatientFullInfo      `json:"patient"`
	Relative Relative     `json:"relative"`
	Officer  Officer     `json:"officer"`
}

type PatientFullInfo struct {
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
	Address     		Address            `json:"address"`
	Weight      		float32            `json:"weight"`
	Height      		float32            `json:"height"`
	Nationality  		string     		   `json:"nationality"`
	Ethnicity    		string       	   `json:"ethnicity"`
	DOB          		string  		   `json:"dob"`          
}

type AppointInfoRes struct {
	Success     bool           		`json:"success"`
	Message     string         		`json:"message"`
	Data        []AppointData   	`json:"data"`
}

type AppointData struct {
	PatientID     uuid.UUID          `json:"patient_id"`
	Fullname      string             `json:"fullname"`
	Hnnumber      string             `json:"hnnumber"`
	Appointments  []AppointDisease   `json:"appointments"`
}

type AppointDisease struct {
	DiseaseID   uuid.UUID              `json:"disease_id"`
	DiseaseName string                 `json:"disease_name"`
	Appointments []AppointmentFullInfo `json:"appointments"`
}

type AppointmentFullInfo struct {
	ID 			uuid.UUID      `json:"id"`
	No 	 		int            `json:"no"`
	Date    	string         `json:"date"`
	StartTime 	string 		   `json:"start_time"`
	EndTime   	string 		   `json:"end_time"`
	Symptom 	string         `json:"symptom"`
	Note    	string         `json:"note"`
	Place   	string         `json:"place"`
	Purpose 	string  	   `json:"purpose"`
	Doctor  	string         `json:"doctor"`
	Officer 	string         `json:"officer"`
	Status  	string         `json:"status"`
	Letter  	bool           `json:"letter"`
	Delay   	bool           `json:"delay"`
	ColorStatus string 		   `json:"color_status"`
}

type VaccineInfoRes struct {
	Success     bool           		`json:"success"`
	Message     string         		`json:"message"`
	Data        []VaccineData   	`json:"data"`
}

type VaccineData struct {
	PatientID     uuid.UUID          `json:"patient_id"`
	Fullname      string             `json:"fullname"`
	Hnnumber      string             `json:"hnnumber"`
	Vaccine       []Vaccine          `json:"vaccine"`
}

type Vaccine struct {
	VaccineID   uuid.UUID `json:"vaccine_id"`
	VaccineName		string    `json:"name"`
	Vaccine     []VaccineFullInfo   `json:"vaccine"`
}

type VaccineFullInfo struct {
	RecordID    uuid.UUID `json:"record_id"`
	Date         string    `json:"date"`
	Type         string    `json:"type"`
	Effect       string    `json:"effect"`
	Note         string    `json:"note"`
	Status       string    `json:"status"`
	Age          string    `json:"age"`
}

type PatientUpdateReq struct {
	Title        string  `json:"title"`
	FirstName    string  `json:"firstname"`
	LastName     string  `json:"lastname"`
	Sex          string  `json:"sex"`
	DOB          string  `json:"dob"`
	Weight       float32 `json:"weight"`
	Height       float32 `json:"height"`
	IDCard       string  `json:"idcard"`
	Rights       string  `json:"rights"`
	Nationality  string  `json:"nationality"`
	Ethnicity    string  `json:"ethnicity"`
	PhoneNumber  string  `json:"phonenumber"`
	Address      Address `json:"address"`
	Diseases     []Disease `json:"diseases"`
}

type PatientUpdateRes struct {
	ID          uuid.UUID `json:"id"`
	Fullname    string    `json:"fullname"`
	HnNumber    string    `json:"hnnumber"`
	IDCard      string    `json:"idcard"`
	PhoneNumber string    `json:"phonenumber"`
	Rights      string    `json:"rights"`
	Nationality string    `json:"nationality"`
	Ethnicity   string    `json:"ethnicity"`
}


type PatientBasicInfoRes struct {
	PatientID uuid.UUID `json:"patient_id"`
	FullName  string    `json:"fullname"`
	HnNumber  string    `json:"hn_number"`
	AgeYears  int       `json:"age_years"`
	AgeMonths int       `json:"age_months"`
	AgeDays   int       `json:"age_days"`
}

type PatientDiseaseQuery struct {
	Type string `query:"type"`
}

type PatientUsecase interface {
	CreateFull(req *PatientFullCreateReq, creatorID *uuid.UUID) (*PatientCreateRes, error)
}

type Transaction interface {
	Do(fn func(repos RepositorySet) error) error
}
type UserRepository interface {
	Create(username, password, role string) (*databases.User, error)
}

type PatientCreateRepo interface {
	CreateWithUser(req *PatientCreateReq, userID uuid.UUID, adminID *uuid.UUID) (*PatientCreateRes, error)
	GenerateNextHnID() (string, error)
}

type PatientGetRepo interface {
	GetPatients(page, limit int) ([]databases.Patient, error)
	GetPatientsWithFilter(req PatientQueryParams) ([]databases.Patient, error)
	CountPatients() (int64, error)
	CountPatientsWithFilter(req PatientQueryParams) (int64, error)
	GetPatientInfoByID(id uuid.UUID) (*databases.Patient, error)
	GetPatientDiseasesInfoByID(id uuid.UUID, diseaseID uuid.UUID) (*databases.Patient, error)
	GetPatientAppointmentsInfoByID(id uuid.UUID) (*databases.Patient, error)
	GetFullNameByUserID(userID uuid.UUID) (string, error)
	GetPatientVaccinesByID(id uuid.UUID) (*databases.Patient , []databases.Vaccine, error)
	GetPatientBasicInfoByID(id uuid.UUID) (*databases.Patient, error)
	GetPatientDiseases(patientID uuid.UUID, dtype string) ([]databases.Disease, error)
}

type PatientGetUsecase interface {
	GetPatientList(page, limit int) (*PatientListRes, error)
	GetPatientListWithFilter(req PatientQueryParams) (*PatientListRes, error)
	GetPatientInfoByID(id uuid.UUID) (*PatientInfoRes, error)
	GetPatientDiseasesInfo(id uuid.UUID, diseaseID uuid.UUID) (*DiseaseInfoRes, error)
	GetPatientAppointmentsInfoByID(id uuid.UUID) (*AppointInfoRes, error)
	GetPatientVaccinesByID(patientID uuid.UUID) (*VaccineInfoRes, error)
	GetPatientBasicInfoByID(id uuid.UUID) (*PatientBasicInfoRes, error)
	GetPatientDiseases(patientID uuid.UUID, dtype string) ([]Disease, error)
}

type PatientUpdateRepo interface {
	UpdatePatientInfo(id uuid.UUID, req *PatientUpdateReq, adminID *uuid.UUID) (*PatientUpdateRes, error)
	UpdatePatientDiseases(patientID uuid.UUID, newDiseases []Disease, adminID *uuid.UUID) error
}

type PatientUpdateUsecase interface {
	UpdatePatientInfo(id uuid.UUID, req *PatientUpdateReq, adminID uuid.UUID) (*PatientUpdateRes, error)
}