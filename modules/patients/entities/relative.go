package entities

type PatientDiseaseEntity struct {
	DiseaseID uint   `json:"disease_id"`
	Name      string `json:"name"`
}

type RelativeEntity struct {
	Title       string     `json:"title"`
	FirstName   string     `json:"firstname"`
	LastName    string     `json:"lastname"`
	PhoneNumber string     `json:"phonenumber"`
	Address     AddressReq `json:"address"`
	Role        string     `json:"role"`
	PatientID   uint       `json:"patient_id"`
	CreatedBy   uint       `json:"created_by"`
	UpdatedBy   uint       `json:"updated_by"`
}

type RelativeDetail struct {
	Title       string     `json:"title"`
	FirstName   string     `json:"firstname"`
	LastName    string     `json:"lastname"`
	PhoneNumber string     `json:"phonenumber"`
	Address     AddressReq `json:"address"`
	CreatedBy   uint       `json:"created_by"`
	UpdatedBy   uint       `json:"updated_by"`
}

type RelativeRepository interface {
	Create(relatives []RelativeEntity) error
}
