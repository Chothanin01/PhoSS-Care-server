package entities

import "github.com/chothanin01/PhoSS-Care-server/pkg/databases"

type AddressReq struct {
	HouseNumber string `json:"house_number"`
	VillageNumber         string `json:"village_number"`
	Alley         string `json:"alley"`
	Road        string `json:"road"`
	SubDistrict string `json:"subdistrict"`
	District    string `json:"district"`
	Province    string `json:"province"`
	ZipCode     string `json:"zipcode"`
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
		DOB          string     `json:"dob"`
		IDCard       string     `json:"idcard"`
		Rights       string     `json:"rights"`
		Nationality  string     `json:"nationality"`
		Ethnicity    string     `json:"ethnicity"`
		PhoneNumber  string     `json:"phonenumber"`
		Address      AddressReq `json:"address"`
		Allergy      string     `json:"allergy"`
		Diseases     []Disease  `json:"diseases"`
	} `json:"patient"`

	Relative struct {
		Kin       RelativeDetail `json:"kin"`
		Caretaker RelativeDetail `json:"caretaker"`
		Medicine  RelativeDetail `json:"medicine"`
	} `json:"relative"`

	Officer struct {
		House RelativeDetail `json:"house"`
		Nurse RelativeDetail `json:"nurse"`
	} `json:"officer"`

	CreatedBy uint `json:"created_by"`
}
type PatientCreateReq struct {
	Title        string     `json:"title"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Dob          string     `json:"dob"`
	HnID         string     `json:"hn_id"`
	IDCard       string     `json:"id_card"`
	Rights       string     `json:"rights"`
	Nationality  string     `json:"nationality"`
	Ethnicity    string     `json:"ethnicity"`
	PhoneNumber  string     `json:"phone_number"`
	Address      AddressReq `json:"address"`
	Allergy      string     `json:"allergy"`
	UserID       uint       `json:"user_id"`
	CreatedBy    uint       `json:"created_by"`
	UpdatedBy    uint       `json:"updated_by"`
}

type PatientCreateRes struct {
	Id          uint64 `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	HnID        string   `json:"hn_id"`
	IDCard      string `json:"id_card"`
	PhoneNumber string `json:"phone_number"`
	Rights      string `json:"rights"`
	Nationality string `json:"nationality"`
	Ethnicity   string `json:"ethnicity"`
}

type RepositorySet struct {
	UserRepo     UserRepository
	PatientRepo  PatientRepository
	RelativeRepo RelativeRepository
	DiseaseRepo  DiseaseRepository
}

type PatientListRes struct {
	Success     bool           `json:"success"`
	Message     string         `json:"message"`
	Page        int            `json:"page"`
	PerPage     int            `json:"per_page"`
	TotalPages  int            `json:"total_pages"`
	Data        []PatientInfo  `json:"data"`
}

type PatientInfo struct {
	ID        uint               `json:"id"`
	FullName  string             `json:"fullname"`
	IDCard    string             `json:"idcard"`
	HnNumber  string             `json:"hnnumber"`
	Diseases  []DiseaseWithStatus `json:"diseases"`
}

type PatientQueryParams struct {
	Search   string   `query:"search"`
	Diseases []string `query:"disease"`
	Appoint  *bool    `query:"appoint"`
	Page     int      `query:"page"`
	Limit    int      `query:"limit"`
}

type PatientUsecase interface {
	CreateFull(req *PatientFullCreateReq) (*PatientCreateRes, error)
}

type Transaction interface {
	Do(fn func(repos RepositorySet) error) error
}
type UserRepository interface {
	Create(username, password, role string) (*databases.User, error)
}

type PatientRepository interface {
	CreateWithUser(req *PatientCreateReq, userID uint) (*PatientCreateRes, error)
	GenerateNextHnID() (string, error)
}

type PatientGetRepo interface {
	GetPatients(page, limit int) ([]databases.Patient, error)
	GetPatientsWithFilter(req PatientQueryParams) ([]databases.Patient, error)
	CountPatients() (int64, error)
	CountPatientsWithFilter(req PatientQueryParams) (int64, error)
}

type PatientGetUsecase interface {
	GetPatientList(page, limit int) (*PatientListRes, error)
	GetPatientListWithFilter(req PatientQueryParams) (*PatientListRes, error)
}

