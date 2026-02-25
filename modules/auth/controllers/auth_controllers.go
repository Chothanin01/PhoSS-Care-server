package controllers

import (
	"github.com/chothanin01/PhoSS-Care-server/modules/auth/entities"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	usecase entities.AuthUsecase
}

func NewAuthController(r fiber.Router, usecase entities.AuthUsecase) {
	controller := &AuthController{
		usecase: usecase,
	}
	r.Post("/login", controller.Login)
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req entities.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	res, err := c.usecase.Login(&req)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(res)
}