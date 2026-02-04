package databases

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string    `gorm:"size:255;not null;unique" json:"username"`
	Password string    `gorm:"size:255;not null" json:"password"`
	Role     string    `gorm:"size:50;not null" json:"role"`
	Admins   []Admin   `json:"admins"`
	Patients []Patient `json:"patients"`
}

type Admin struct {
	gorm.Model
	Title		   string `gorm:"size:50;not null" json:"title"`
	FirstName      string `gorm:"size:255;not null" json:"first_name"`
	LastName       string `gorm:"size:255;not null" json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	UserID         uint   `json:"user_id"`
	User           User   `gorm:"foreignKey:UserID" json:"user"`
	CreatedBy      uint
	CreatedByUser  User `gorm:"foreignKey:CreatedBy" json:"created_by_user"`
	UpdatedBy      uint
	UpdatedByUser  User `gorm:"foreignKey:UpdatedBy" json:"updated_by_user"`
}

type Patient struct {
	gorm.Model
	Title		   string 	 `gorm:"size:50;not null" json:"title"`
	FirstName      string    `gorm:"size:255;not null" json:"first_name"`
	LastName       string    `gorm:"size:255;not null" json:"last_name"`
	DOB            time.Time `json:"dob"`
	HnID           string    `gorm:"uniqueIndex;size:7;not null" json:"hn_id"`
	IDCard         string    `gorm:"size:13;uniqueIndex;not null" json:"id_card"`
	Rights         string    `gorm:"size:255;not null" json:"rights"`
	Nationality    string    `gorm:"size:255;not null" json:"nationality"`
	Ethnicity      string    `gorm:"size:255;not null" json:"ethnicity"`
	PhoneNumber    string    `gorm:"size:10" json:"phone_number"`
	Address        Address   `gorm:"type:json" json:"address"`
	Allergy        string    `json:"allergy"`
	ProfilePicture string    `json:"profile_picture"`
	UserID         uint      `json:"user_id"`
	User           User      `gorm:"foreignKey:UserID" json:"user"`
	CreatedBy      uint
	CreatedByUser  User `gorm:"foreignKey:CreatedBy" json:"created_by_user"`
	UpdatedBy      uint
	UpdatedByUser  User `gorm:"foreignKey:UpdatedBy" json:"updated_by_user"`
	Healths        []Health          `json:"healths"`
	Relatives      []Relative        `json:"relatives"`
	Appointments   []Appoint         `json:"appointments"`
	Diseases       []PatientDisease  `json:"diseases"`
}

type Health struct {
	gorm.Model
	Weight        float64 `gorm:"type:decimal(5,2);not null" json:"weight"`
	Height        int     `gorm:"not null" json:"height"`
	BMI           float64 `gorm:"type:decimal(5,2);not null" json:"bmi"`
	Pulse         int     `gorm:"not null" json:"pulse"`
	Sugar         int     `json:"sugar"`
	PatientID     uint    `json:"patient_id"`
	Patient       Patient `gorm:"foreignKey:PatientID" json:"patient"`
	AppointID     uint    `json:"appoint_id"`
	Appoint       Appoint `gorm:"foreignKey:AppointID" json:"appoint"`
	CreatedBy     uint
	CreatedByUser User `gorm:"foreignKey:CreatedBy" json:"created_by_user"`
	UpdatedBy     uint
	UpdatedByUser User `gorm:"foreignKey:UpdatedBy" json:"updated_by_user"`
}

type Relative struct {
	gorm.Model
	Title		  string   `gorm:"size:50;not null" json:"title"`
	FirstName     string   `gorm:"size:255;not null" json:"first_name"`
	LastName      string   `gorm:"size:255;not null" json:"last_name"`
	PhoneNumber   string   `gorm:"size:10" json:"phone_number"`
	Address       Address  `gorm:"type:json" json:"address"`
	Role          string   `gorm:"size:255;not null" json:"role"`
	PatientID     uint     `json:"patient_id"`
	Patient       Patient  `gorm:"foreignKey:PatientID" json:"patient"`
	CreatedBy     uint
	CreatedByUser User `gorm:"foreignKey:CreatedBy" json:"created_by_user"`
	UpdatedBy     uint
	UpdatedByUser User `gorm:"foreignKey:UpdatedBy" json:"updated_by_user"`
}

type Appoint struct {
	gorm.Model
	No           int       `json:"no"`
	Date         time.Time `json:"date"`
	Time         string    `gorm:"size:13" json:"time"`
	Symptom      string    `json:"symptom"`
	Note         string    `json:"note"`
	Place        string    `json:"place"`
	Doctor       string    `json:"doctor"`
	Status       string    `json:"status"`
	Letter       bool      `json:"letter"`
	Delay        bool      `json:"delay"`
	PatientID    uint      `json:"patient_id"`
	Patient      Patient   `gorm:"foreignKey:PatientID" json:"patient"`
	DiseaseID    uint      `json:"disease_id"`
	Disease      Disease   `gorm:"foreignKey:DiseaseID" json:"disease"`
	CreatedBy    uint
	CreatedByUser  User `gorm:"foreignKey:CreatedBy" json:"created_by_user"`
	UpdatedBy    uint
	UpdatedByUser  User `gorm:"foreignKey:UpdatedBy" json:"updated_by_user"`
	Healths      []Health           `json:"healths"`
	Vaccinations []VaccinationRecord `json:"vaccinations"`
	Requests     []Request          `json:"requests"`
}

type Disease struct {
	gorm.Model
	Name     string            `gorm:"size:255;not null" json:"name"`
	Patients []PatientDisease  `json:"patients"`
}

type PatientDisease struct {
	gorm.Model
	PatientID     uint      `json:"patient_id"`
	Patient       Patient   `gorm:"foreignKey:PatientID" json:"patient"`
	DiseaseID     uint      `json:"disease_id"`
	Disease       Disease   `gorm:"foreignKey:DiseaseID" json:"disease"`
	CreatedBy     uint
	CreatedByUser User `gorm:"foreignKey:CreatedBy" json:"created_by_user"`
	UpdatedBy     uint
	UpdatedByUser User `gorm:"foreignKey:UpdatedBy" json:"updated_by_user"`
}

type Request struct {
	gorm.Model
	RequestType   string        `gorm:"size:255;not null" json:"request_type"`
	Description   string        `json:"description"`
	Date          time.Time     `json:"date"`
	Time          string        `gorm:"size:13" json:"time"`
	Status        string        `gorm:"size:50;not null" json:"status"`
	PatientID     uint          `json:"patient_id"`
	Patient       Patient       `gorm:"foreignKey:PatientID" json:"patient"`
	AppointID     uint          `json:"appoint_id"`
	Appoint       Appoint       `gorm:"foreignKey:AppointID" json:"appoint"`
	UpdatedBy     uint
	UpdatedByUser User          `gorm:"foreignKey:UpdatedBy" json:"updated_by_user"`
	Notifications []Notification `json:"notifications"`
}

type Vaccine struct {
	gorm.Model
	Name    string              `gorm:"not null" json:"name"`
	Age     string              `gorm:"not null" json:"age"`
	Effect  string              `gorm:"not null" json:"effect"`
	Note    string              `json:"note"`
	Records []VaccinationRecord `json:"records"`
}

type VaccinationRecord struct {
	gorm.Model
	VaccineID     uint      `json:"vaccine_id"`
	Vaccine       Vaccine   `gorm:"foreignKey:VaccineID" json:"vaccine"`
	AppointID     uint      `json:"appoint_id"`
	DoseNumber	  int       `json:"dose_number"`
	Appoint       Appoint   `gorm:"foreignKey:AppointID" json:"appoint"`
	CreatedBy     uint
	CreatedByUser User      `gorm:"foreignKey:CreatedBy" json:"created_by_user"`
	UpdatedBy     uint
	UpdatedByUser User      `gorm:"foreignKey:UpdatedBy" json:"updated_by_user"`
}

type Notification struct {
	gorm.Model
	Name          string `gorm:"size:255;not null" json:"name"`
	Status        string `gorm:"size:50;not null" json:"status"`
	Note          string `json:"note"`
	RequestID     uint   `json:"request_id"`
	Request       Request `gorm:"foreignKey:RequestID" json:"request"`
	CreatedBy     uint
	CreatedByUser User   `gorm:"foreignKey:CreatedBy" json:"created_by_user"`
	UpdatedBy     uint
	UpdatedByUser User   `gorm:"foreignKey:UpdatedBy" json:"updated_by_user"`
}

type Address struct {
	HouseNumber 		  string `json:"house_number"`
	VillageNumber         string `json:"village_number"`
	Alley         		  string `json:"alley"`
	Road        		  string `json:"road"`
	SubDistrict 		  string `json:"subdistrict"`
	District    		  string `json:"district"`
	Province    		  string `json:"province"`
	ZipCode     	      string `json:"zipcode"`
}


