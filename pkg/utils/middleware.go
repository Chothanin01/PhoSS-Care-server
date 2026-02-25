package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func NewJWTMiddleware(jwtSvc JWTService, allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization header"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}

		token, err := jwtSvc.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		role := claims["role"]
		if len(allowedRoles) > 0 {
			authorized := false
			for _, r := range allowedRoles {
				if r == role {
					authorized = true
					break
				}
			}
			if !authorized {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied for this role"})
			}
		}

		c.Locals("user", claims)
		return c.Next()
	}
}