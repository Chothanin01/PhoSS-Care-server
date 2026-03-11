package databases

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type User struct {
	BaseModel
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
	Role     string `gorm:"size:50;not null" json:"role"`
}

func (User) TableName() string { return "users" }


type Admin struct {
	BaseModel
	Title          string    `gorm:"size:50;not null" json:"title"`
	FirstName      string    `gorm:"size:255;not null" json:"first_name"`
	LastName       string    `gorm:"size:255;not null" json:"last_name"`
	ProfilePicture string    `json:"profile_picture"`
	UserID         uuid.UUID `json:"user_id"`
	User           User      `gorm:"foreignKey:UserID"`

	CreatedBy     *uuid.UUID
	UpdatedBy     *uuid.UUID
	CreatedByUser *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Patient struct {
	BaseModel
	Title       string    `gorm:"size:50;not null" json:"title"`
	FirstName   string    `gorm:"size:255;not null" json:"first_name"`
	LastName    string    `gorm:"size:255;not null" json:"last_name"`
	Sex         string    `gorm:"size:10;not null" json:"sex"`
	DOB         time.Time `json:"dob"`
	HnID        string    `gorm:"uniqueIndex;size:7;not null" json:"hn_id"`
	IDCard      string    `gorm:"size:13;uniqueIndex;not null" json:"id_card"`
	Rights      string    `gorm:"size:255;not null" json:"rights"`
	Nationality string    `gorm:"size:255;not null" json:"nationality"`
	Ethnicity   string    `gorm:"size:255;not null" json:"ethnicity"`
	PhoneNumber string    `gorm:"size:10" json:"phone_number"`
	Address     Address   `gorm:"type:jsonb" json:"address"`
	Allergy     string    `json:"allergy"`
	Weight      float32   `gorm:"not null" json:"weight"`
	Height      float32   `gorm:"not null" json:"height"`

	UserID uuid.UUID `json:"user_id"`
	User   User      `gorm:"foreignKey:UserID"`

	Relatives    []Relative       `json:"relatives"`
	Appointments []Appoint        `gorm:"foreignKey:PatientID;references:ID" json:"appointments"`
	Diseases     []PatientDisease `json:"diseases"`
	Healths      []Health         `json:"healths"`

	CreatedBy     *uuid.UUID
	UpdatedBy     *uuid.UUID
	CreatedByUser *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Health struct {
	BaseModel
	Weight    float64   `gorm:"type:decimal(5,2);not null" json:"weight"`
	Height    int       `gorm:"not null" json:"height"`
	BMI       float64   `gorm:"type:decimal(5,2);not null" json:"bmi"`
	Pulse     int       `json:"pulse"`
	Sugar     int       `json:"sugar"`
	PatientID uuid.UUID `json:"patient_id"`
	AppointID uuid.UUID `json:"appoint_id"`

	Patient Patient `gorm:"foreignKey:PatientID"`
	Appoint Appoint `gorm:"foreignKey:AppointID"`

	CreatedBy     *uuid.UUID
	UpdatedBy     *uuid.UUID
	CreatedByUser *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Relative struct {
	BaseModel
	Title       string  `gorm:"size:50;not null" json:"title"`
	FirstName   string  `gorm:"size:255;not null" json:"first_name"`
	LastName    string  `gorm:"size:255;not null" json:"last_name"`
	PhoneNumber string  `gorm:"size:10" json:"phone_number"`
	Address     Address `gorm:"type:jsonb" json:"address"`
	Role        string  `gorm:"size:255;not null" json:"role"`
	PatientID   uuid.UUID
	Patient     Patient `gorm:"foreignKey:PatientID"`

	CreatedBy     *uuid.UUID
	UpdatedBy     *uuid.UUID
	CreatedByUser *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Appoint struct {
	BaseModel
	No        int       `json:"no"`
	Date      time.Time `json:"date"`
	StartTime string `gorm:"size:5" json:"start_time"`
	EndTime   string `gorm:"size:5" json:"end_time"`
	Symptom   string    `json:"symptom"`
	Note      string    `json:"note"`
	Place     string    `json:"place"`
	Doctor    string    `json:"doctor"`
	Status    string    `json:"status"`
	Purpose   string    `json:"purpose"`
	Letter    bool      `json:"letter"`
	Delay     bool      `json:"delay"`
	PatientID uuid.UUID `json:"patient_id"`
	DiseaseID uuid.UUID `json:"disease_id"`

	Patient Patient `gorm:"foreignKey:PatientID"`
	Disease Disease `gorm:"foreignKey:DiseaseID"`

	Vaccinations []VaccinationRecord `json:"vaccinations"`
	Requests     []Request           `json:"requests"`
	Healths      []Health            `json:"healths"`

	CreatedBy     *uuid.UUID
	UpdatedBy     *uuid.UUID
	CreatedByUser *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Disease struct {
	BaseModel
	Name     string           `gorm:"size:255;not null" json:"name"`
	Patients []PatientDisease `json:"patients"`

	CreatedBy     *uuid.UUID
	UpdatedBy     *uuid.UUID
	CreatedByUser *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type PatientDisease struct {
	BaseModel
	PatientID uuid.UUID
	DiseaseID uuid.UUID
	Patient   *Patient `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Disease   *Disease `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedBy     *uuid.UUID
	UpdatedBy     *uuid.UUID
	CreatedByUser *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (PatientDisease) TableName() string { return "patient_disease" }

type Vaccine struct {
	BaseModel
	Name   string `gorm:"not null" json:"name"`
	Age    string `gorm:"not null" json:"age"`
	Type   string `gorm:"not null" json:"type"`
	Effect string `gorm:"not null" json:"effect"`
	Note   string `json:"note"`

	Records []VaccinationRecord `json:"records"`

	CreatedBy     *uuid.UUID
	UpdatedBy     *uuid.UUID
	CreatedByUser *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type VaccinationRecord struct {
	BaseModel
	VaccineID  uuid.UUID
	AppointID  uuid.UUID
	DoseNumber int
	VaccineDoctor string `gorm:"size:50;not null"`
	Status     string `gorm:"size:50;not null"`

	Vaccine Vaccine `gorm:"foreignKey:VaccineID"`
	Appoint Appoint `gorm:"foreignKey:AppointID"`

	CreatedBy     *uuid.UUID
	UpdatedBy     *uuid.UUID
	CreatedByUser *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Request struct {
	BaseModel
	RequestType string    `gorm:"size:255;not null"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	StartTime string `gorm:"size:5" json:"start_time"`
	EndTime   string `gorm:"size:5" json:"end_time"`
	Status      string    `gorm:"size:50;not null"`
	PatientID   uuid.UUID
	AppointID   uuid.UUID
	DiseaseID uuid.UUID   `json:"disease_id"`
	
	Patient Patient `gorm:"foreignKey:PatientID"`
	Appoint Appoint `gorm:"foreignKey:AppointID"`
	Disease Disease `gorm:"foreignKey:DiseaseID"`

	Notifications []Notification `json:"notifications"`

	CreatedBy     *uuid.UUID
	UpdatedBy     *uuid.UUID
	CreatedByUser *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Notification struct {
	BaseModel
	Name      string    `gorm:"size:255;not null"`
	Status    string    `gorm:"size:50;not null"`
	Note      string    `json:"note"`
	RequestID uuid.UUID `json:"request_id"`

	Request Request `gorm:"foreignKey:RequestID"`

	CreatedBy     *uuid.UUID
	UpdatedBy     *uuid.UUID
	CreatedByUser *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
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

func (a Address) Value() (driver.Value, error) { return json.Marshal(a) }
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

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc)
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

func (b *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	b.UpdatedAt = time.Now().In(loc)
	return nil
}

