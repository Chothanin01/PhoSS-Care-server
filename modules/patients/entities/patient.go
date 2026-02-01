package entities


type PatientUsecase interface {
	Create(req *PatientCreateReq) (*PatientCreateRes, error)
}

type PatientRepository interface {
	Create(req *PatientCreateReq) (*PatientCreateRes, error)
}

type PatientCreateReq struct {
	FirstName      string    `json:"firstName" db:"first_name"`
	LastName       string    `json:"lastName" db:"last_name"`
	Dob            string `json:"dob" db:"dob"`
	HnID           uint      `json:"hnId" db:"hn_id"`
	IDCard         string    `json:"idCard" db:"id_card"`
	Rights         string    `json:"rights" db:"rights"`
	Nationality    string    `json:"nationality" db:"nationality"`
	Ethnicity      string    `json:"ethnicity" db:"ethnicity"`
	PhoneNumber    string    `json:"phoneNumber" db:"phone_number"`
	Address        string    `json:"address" db:"address"`
	Allergy        string    `json:"allergy" db:"allergy"`
	ProfilePicture string    `json:"profilePicture" db:"profile_picture"`
	UserID         uint      `json:"userId" db:"user_id"`
}

type PatientCreateRes struct {
	Id          uint64 `json:"id" db:"id"`
	FirstName   string `json:"firstName" db:"first_name"`
	LastName    string `json:"lastName" db:"last_name"`
	HnID        uint   `json:"hnId" db:"hn_id"`
	IDCard      string `json:"idCard" db:"id_card"`
	PhoneNumber string `json:"phoneNumber" db:"phone_number"`
	Rights      string `json:"rights" db:"rights"`
	Nationality string `json:"nationality" db:"nationality"`
	Ethnicity   string `json:"ethnicity" db:"ethnicity"`
}
