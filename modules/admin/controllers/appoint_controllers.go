package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
)

type AppointmentController struct {
	usecase entities.AppointmentUsecase
}

func NewAppointmentController(r fiber.Router, uc entities.AppointmentUsecase) {
	controller := &AppointmentController{usecase: uc}
	r.Post("/", controller.CreateAppointment)
	r.Get("/:id/vaccination", controller.FindOngoingVaccination)
}

func (c *AppointmentController) CreateAppointment(ctx *fiber.Ctx) error {
	var req entities.AppointmentCreateReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	claims := ctx.Locals("user").(jwt.MapClaims)
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
	adminID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user_id"})
	}

	res, err := c.usecase.CreateAppointment(&req, adminID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Appointment created successfully",
		"data":    res,
	})
}

func (c *AppointmentController) FindOngoingVaccination(ctx *fiber.Ctx) error {

	idParam := ctx.Params("id")

	patientID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid patient id",
		})
	}

	data, err := c.usecase.FindOngoingVaccination(patientID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "no ongoing vaccination",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Get ongoing vaccination successfully.",
		"data": data,
	})
}