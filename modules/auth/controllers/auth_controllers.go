package controllers

import (
	"time"

	"github.com/chothanin01/PhoSS-Care-server/modules/auth/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	usecase    entities.AuthUsecase
	jwtService utils.JWTService
}

func NewAuthController(r fiber.Router, middleware fiber.Handler, usecase entities.AuthUsecase, jwt utils.JWTService) {
	controller := &AuthController{
		usecase:    usecase,
		jwtService: jwt,
	}
	
	r.Post("/login", controller.Login)
	r.Get("/logout", middleware, controller.Logout)

	r.Get("/refresh", middleware, controller.RefreshToken)
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

func (c *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid token structure"})
	}
	
	userID, _ := claims["user_id"].(string)
	roleID, _ := claims["role_id"].(string)
	role, _ := claims["role"].(string)

	newToken, err := c.jwtService.GenerateToken(userID, role, roleID, nil)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate new token",
		})
	}

	ctx.Cookie(&fiber.Cookie{
    Name:     "accessToken",
    Value:    newToken,
    Expires:  time.Now().Add(time.Hour * 24),
    HTTPOnly: true, 
    Secure:   true,
    SameSite: "Lax",
	})

	return ctx.JSON(fiber.Map{
		"token": newToken,
	})
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {

	ctx.Cookie(&fiber.Cookie{
		Name:     "token", 
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), 
		HTTPOnly: true,
		Secure:   true, 
	})

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Successfully logged out",
	})
}