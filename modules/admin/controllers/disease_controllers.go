package controllers

import (

	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
	"github.com/gofiber/fiber/v2"
)

type DiseaseController struct {
	DiseaseUsecase entities.DiseaseUsecase
}

func NewDiseaseController(r fiber.Router, diseaseUC entities.DiseaseUsecase) {
	controller := &DiseaseController{
		DiseaseUsecase: diseaseUC,
	}
	r.Get("/", controller.GetAllDiseases)
	r.Get("vaccines", controller.GetAllVaccines)
}

func (c *DiseaseController) GetAllDiseases(ctx *fiber.Ctx) error {
	diseases, err := c.DiseaseUsecase.GetAllDiseases()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.GetDiseaseListRes{
		Diseases: diseases,
	})
}

func (c *DiseaseController) GetAllVaccines(ctx *fiber.Ctx) error {

	vaccines, err := c.DiseaseUsecase.GetAllVaccines()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	message := "Vaccine list retrieved successfully"
	if len(vaccines) == 0 {
		message = "No vaccine found"
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    vaccines,
	})
}