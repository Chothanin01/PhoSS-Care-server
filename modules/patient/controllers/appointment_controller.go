package controllers

import (

	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AppointmentController struct {
	GetUsecase entities.GetAppointmentUsecase
}

func NewAppointmentController(r fiber.Router, uc entities.GetAppointmentUsecase) {
	controller := &AppointmentController{GetUsecase: uc}
	r.Get("/:disease_id", controller.GetAppointmentDetail)
}

func (c *AppointmentController) GetAppointmentDetail(ctx *fiber.Ctx) error {

	claims := ctx.Locals("user").(jwt.MapClaims)
	patientIDStr := claims["role_id"].(string)
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user id"})
	}

	diseaseIDStr := ctx.Params("disease_id")
	diseaseID, err := uuid.Parse(diseaseIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid disease id"})
	}

	data, err := c.GetUsecase.GetAppointmentDetail(patientID, diseaseID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "appointment not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}
