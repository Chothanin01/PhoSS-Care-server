package databases

import (
	"time"
	"encoding/json"
	"fmt"

	"database/sql/driver"
	"database/sql"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:false" json:"updated_at"`
	CreatedBy *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID     `gorm:"type:uuid" json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type User struct {
	BaseModel
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
	Role     string `gorm:"size:50;not null" json:"role"`
	Admin    []Admin  `gorm:"-:migration" json:"admin"`
	Patient  []Patient `gorm:"-:migration" json:"patients"`
}

func (User) TableName() string {
	return "users"
}

type Admin struct {
	BaseModel
	Title           string    `gorm:"size:50;not null" json:"title"`
	FirstName       string    `gorm:"size:255;not null" json:"first_name"`
	LastName        string    `gorm:"size:255;not null" json:"last_name"`
	ProfilePicture  string    `json:"profile_picture"`
	UserID          uuid.UUID `json:"user_id"`
	User            User      `gorm:"foreignKey:UserID"`
	CreatedBy       uuid.UUID
	CreatedByUser   User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
	UpdatedBy       uuid.UUID
	UpdatedByUser   User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
}

type Patient struct {
	BaseModel
	Title           string    `gorm:"size:50;not null" json:"title"`
	FirstName       string    `gorm:"size:255;not null" json:"first_name"`
	LastName        string    `gorm:"size:255;not null" json:"last_name"`
	Sex 		    string    `gorm:"size:10;not null" json:"sex"`
	DOB             time.Time `json:"dob"`
	HnID            string    `gorm:"uniqueIndex;size:7;not null" json:"hn_id"`
	IDCard          string    `gorm:"size:13;uniqueIndex;not null" json:"id_card"`
	Rights          string    `gorm:"size:255;not null" json:"rights"`
	Nationality     string    `gorm:"size:255;not null" json:"nationality"`
	Ethnicity       string    `gorm:"size:255;not null" json:"ethnicity"`
	PhoneNumber     string    `gorm:"size:10" json:"phone_number"`
	Address         Address   `gorm:"type:jsonb" json:"address"`
	Allergy         string    `json:"allergy"`
	ProfilePicture  string    `json:"profile_picture"`
	UserID          uuid.UUID `json:"user_id"`
	User            User      `gorm:"foreignKey:UserID"`
	CreatedBy       uuid.UUID
	CreatedByUser   User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
	UpdatedBy       uuid.UUID
	UpdatedByUser   User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
	Healths         []Health         `json:"healths"`
	Relatives       []Relative       `json:"relatives"`
	Appointments    []Appoint        `gorm:"foreignKey:PatientID;references:ID" json:"appointments"`
	Diseases        []PatientDisease `json:"diseases"`
}

type Health struct {
	BaseModel
	Weight        float64   `gorm:"type:decimal(5,2);not null" json:"weight"`
	Height        int       `gorm:"not null" json:"height"`
	BMI           float64   `gorm:"type:decimal(5,2);not null" json:"bmi"`
	Pulse         int       `gorm:"not null" json:"pulse"`
	Sugar         int       `json:"sugar"`
	PatientID     uuid.UUID `json:"patient_id"`
	Patient       Patient   `gorm:"foreignKey:PatientID"`
	AppointID     uuid.UUID `json:"appoint_id"`
	Appoint       Appoint   `gorm:"foreignKey:AppointID"`
	CreatedBy     uuid.UUID
	CreatedByUser User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
	UpdatedBy     uuid.UUID
	UpdatedByUser User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
}

type Relative struct {
	BaseModel
	Title          string    `gorm:"size:50;not null" json:"title"`
	FirstName      string    `gorm:"size:255;not null" json:"first_name"`
	LastName       string    `gorm:"size:255;not null" json:"last_name"`
	PhoneNumber    string    `gorm:"size:10" json:"phone_number"`
	Address        Address   `gorm:"type:json" json:"address"`
	Role           string    `gorm:"size:255;not null" json:"role"`
	PatientID      uuid.UUID `json:"patient_id"`
	Patient        Patient   `gorm:"foreignKey:PatientID"`
	CreatedBy      uuid.UUID
	CreatedByUser  User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
	UpdatedBy      uuid.UUID
	UpdatedByUser  User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
}

type Appoint struct {
	BaseModel
	No             int       `json:"no"`
	Date           time.Time `json:"date"`
	Time           string    `gorm:"size:13" json:"time"`
	Symptom        string    `json:"symptom"`
	Note           string    `json:"note"`
	Place          string    `json:"place"`
	Doctor         string    `json:"doctor"`
	Status         string    `json:"status"`
	Letter         bool      `json:"letter"`
	Delay          bool      `json:"delay"`
	PatientID      uuid.UUID `json:"patient_id"`
	Patient        Patient   `gorm:"foreignKey:PatientID;references:ID"`
	DiseaseID      uuid.UUID `json:"disease_id"`
	Disease        Disease   `gorm:"foreignKey:DiseaseID"`
	CreatedBy      uuid.UUID
	CreatedByUser  User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
	UpdatedBy      uuid.UUID
	UpdatedByUser  User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
	Healths        []Health           `json:"healths"`
	Vaccinations   []VaccinationRecord `json:"vaccinations"`
	Requests       []Request          `json:"requests"`
}

type Disease struct {
	BaseModel
	Name     string           `gorm:"size:255;not null" json:"name"`
	Patients []PatientDisease `json:"patients"`
}

type PatientDisease struct {
	BaseModel
	PatientID uuid.UUID `json:"patient_id"`
	Patient   *Patient  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"patient,omitempty"`
	DiseaseID uuid.UUID `json:"disease_id"`
	Disease   *Disease  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;->" json:"disease,omitempty"`
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
}

func (PatientDisease) TableName() string { return "patient_disease" }


type Request struct {
	BaseModel
	RequestType   string        `gorm:"size:255;not null" json:"request_type"`
	Description   string        `json:"description"`
	Date          time.Time     `json:"date"`
	Time          string        `gorm:"size:13" json:"time"`
	Status        string        `gorm:"size:50;not null" json:"status"`
	PatientID     uuid.UUID     `json:"patient_id"`
	Patient       Patient       `gorm:"foreignKey:PatientID"`
	AppointID     uuid.UUID     `json:"appoint_id"`
	Appoint       Appoint       `gorm:"foreignKey:AppointID"`
	UpdatedBy     uuid.UUID
	UpdatedByUser User           `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
	Notifications []Notification `json:"notifications"`
}

type Vaccine struct {
	BaseModel
	Name    string              `gorm:"not null" json:"name"`
	Age     string              `gorm:"not null" json:"age"`
	Effect  string              `gorm:"not null" json:"effect"`
	Note    string              `json:"note"`
	Records []VaccinationRecord `json:"records"`
}

type VaccinationRecord struct {
	BaseModel
	VaccineID     uuid.UUID `json:"vaccine_id"`
	Vaccine       Vaccine   `gorm:"foreignKey:VaccineID"`
	AppointID     uuid.UUID `json:"appoint_id"`
	DoseNumber    int       `json:"dose_number"`
	Appoint       Appoint   `gorm:"foreignKey:AppointID"`
	CreatedBy     uuid.UUID
	CreatedByUser User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
	UpdatedBy     uuid.UUID
	UpdatedByUser User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
}

type Notification struct {
	BaseModel
	Name          string    `gorm:"size:255;not null" json:"name"`
	Status        string    `gorm:"size:50;not null" json:"status"`
	Note          string    `json:"note"`
	RequestID     uuid.UUID `json:"request_id"`
	Request       Request   `gorm:"foreignKey:RequestID"`
	CreatedBy     uuid.UUID
	CreatedByUser User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
	UpdatedBy     uuid.UUID
	UpdatedByUser User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;"`
}

type Address struct {
	HouseNumber   string `json:"house_number"`
	VillageNumber string `json:"village_number"`
	Alley         string `json:"alley"`
	Road          string `json:"road"`
	SubDistrict   string `json:"subdistrict"`
	District      string `json:"district"`
	Province      string `json:"province"`
	ZipCode       string `json:"zipcode"`
}

func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Address) Scan(value interface{}) error {
	if value == nil {
		*a = Address{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	default:
		return fmt.Errorf("unsupported type %T for Address", value)
	}
}

var (
	_ driver.Valuer = (*Address)(nil)
	_ sql.Scanner   = (*Address)(nil)
)

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc)
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

func (b *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	b.UpdatedAt = time.Now().In(loc)
	return nil
}