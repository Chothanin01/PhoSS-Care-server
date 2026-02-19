package controllers

import (
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
	"github.com/gofiber/fiber/v2"
)

type RelativeController struct {
	RelativeUsecase entities.RelativeUsecase
}

func NewRelativeController(r fiber.Router, relativeUC entities.RelativeUsecase) {
	controller := &RelativeController{
		RelativeUsecase: relativeUC,
	}

	r.Patch("/:id/relatives", controller.UpdateAllRelatives)
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

	res, err := c.RelativeUsecase.UpdateAllRelatives(patientID, &req)
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
