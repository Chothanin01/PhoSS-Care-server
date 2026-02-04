package controllers

import (
	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
	"github.com/gofiber/fiber/v2"
)

type PatientController struct {
	PatientUsecase entities.PatientUsecase
}

func NewPatientController(r fiber.Router, patientUsecase entities.PatientUsecase) {
	controllers := &PatientController{
		PatientUsecase: patientUsecase,
	}
	r.Post("/", controllers.CreatePatient)
}


func (c *PatientController) CreatePatient(ctx *fiber.Ctx) error {
	req := new(entities.PatientFullCreateReq)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	res, err := c.PatientUsecase.CreateFull(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"message":    "User, Patient, and Relatives created successfully",
		"data":       res,
		"statusCode": fiber.StatusOK,
	})
}
