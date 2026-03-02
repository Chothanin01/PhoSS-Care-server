package controllers

import (
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
	"github.com/gofiber/fiber/v2"
)

type AdminController struct {
	usecase entities.AdminUsecase
}

func NewAdminController(r fiber.Router, uc entities.AdminUsecase) {
	controller := &AdminController{usecase: uc}
	r.Post("/", controller.CreateAdmin)
}

func (c *AdminController) CreateAdmin(ctx *fiber.Ctx) error {
	req := new(entities.AdminCreateReq)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	res, err := c.usecase.CreateAdmin(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Admin created successfully",
		"data":    res,
	})
}