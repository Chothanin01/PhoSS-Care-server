package entities

type Disease struct {
	DiseaseID uint   `json:"disease_id"`
	Name      string `json:"name"`
}

type DiseaseRepository interface {
	LinkPatientDiseases(patientID uint, diseases []PatientDiseaseEntity) error
}

