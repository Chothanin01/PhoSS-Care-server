package controllers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
)

type PatientController struct {
	GetUsecase entities.PatientGetUsecase
}

func NewPatientController(r fiber.Router, uc entities.PatientGetUsecase) {
	controller := &PatientController{GetUsecase: uc}
	r.Get("/basicinfo", controller.GetPatientBasicInfo)
}

func (c *PatientController) GetPatientBasicInfo(ctx *fiber.Ctx) error {

	claims := ctx.Locals("user").(jwt.MapClaims)

	userIDStr := claims["user_id"].(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid user id",
		})
	}

	data, err := c.GetUsecase.GetPatientBasicInfo(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": data,
	})
}
