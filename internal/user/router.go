package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/katsumi/inventory_api/config"
	"github.com/katsumi/inventory_api/internal/auth"
)

func CreateUserGroup(app *fiber.App, userController *UserController, config config.EnvVars) {
	userGroup := app.Group("/users")

	// middleware to protect routes
	authMiddleware := auth.NewAuthMiddleware(config)
	userGroup.Use(authMiddleware.ValidateToken)

	// auth routes
	userGroup.Get("/me", userController.profile)
}
