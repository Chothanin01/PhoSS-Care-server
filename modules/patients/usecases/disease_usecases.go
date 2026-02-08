package usecases

import (

	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
)

// ---------------------- DISEASE ----------------------
type newDiseaseUsecase struct {
	tx entities.Transaction
}

func NewDiseaseUsecase(tx entities.Transaction) *newDiseaseUsecase {
	return &newDiseaseUsecase{
		tx: tx,
	}
}

func (u *newDiseaseUsecase) GetAllDiseases() ([]entities.Disease, error) {
	var res []entities.Disease
	err := u.tx.Do(func(r entities.RepositorySet) error {
		diseases, err := r.DiseaseGetRepo.GetAllDiseases()
		if err != nil {
			return err
		}
		res = diseases
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}