package controllers

import (
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type RelativeController struct {
	RelativeUsecase entities.RelativeUsecase
}

func NewRelativeController(r fiber.Router, relativeUC entities.RelativeUsecase) {
	controller := &RelativeController{
		RelativeUsecase: relativeUC,
	}

	r.Patch("/:id/relatives", controller.UpdateAllRelatives)
	r.Patch("/:id/officer", controller.UpdateOfficers)
}

func (c *RelativeController) UpdateAllRelatives(ctx *fiber.Ctx) error {
	patientIDParam := ctx.Params("id")
	if patientIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Missing patient ID",
		})
	}

	patientID, err := uuid.Parse(patientIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid patient ID format",
		})
	}

	var req entities.RelativeAllUpdateReq
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

	res, err := c.RelativeUsecase.UpdateAllRelatives(patientID, &req, adminID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Relatives updated successfully",
		"data":    res,
	})
}

func (c *RelativeController) UpdateOfficers(ctx *fiber.Ctx) error {
	patientID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid patient ID",
		})
	}

	var req entities.OfficerAllUpdateReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "invalid request body"})
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

	res, err := c.RelativeUsecase.UpdateAllOfficers(patientID, &req, adminID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Officers updated successfully", "data": res})
}