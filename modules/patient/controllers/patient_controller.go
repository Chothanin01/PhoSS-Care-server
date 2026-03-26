package controllers

import (

	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type PatientController struct {
	GetUsecase entities.GetPatientUsecase
}

func NewPatientController(r fiber.Router, uc entities.GetPatientUsecase) {
	controller := &PatientController{GetUsecase: uc}
	r.Get("/basicinfo", controller.GetPatientBasicInfo)
	r.Get("/appointment", controller.GetPatientAppointment)
}

func (c *PatientController) GetPatientBasicInfo(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized access",
		})
	}

	patientIDStr, ok := claims["role_id"].(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid token claims",
		})
	}

	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid user id",
		})
	}

	data, err := c.GetUsecase.GetPatientBasicInfo(patientID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func (c *PatientController) GetPatientAppointment(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized access",
		})
	}

	patientIDStr, ok := claims["role_id"].(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid token claims",
		})
	}

	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid user id",
		})
	}

	data, err := c.GetUsecase.GetPatientAppointment(patientID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}
