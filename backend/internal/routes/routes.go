package routes

import (
	"github.com/ashmitsharp/msg-agg/internal/handlers"
	"github.com/ashmitsharp/msg-agg/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler, identityHandler *handlers.IdentityHandler) {
	api := app.Group("/api")

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/signup", authHandler.SignUp)
	auth.Post("/login", authHandler.Login)

	// Protected routes
	protected := api.Group("/", middlewares.AuthMiddleware())

	// Identity Routes
	identity := protected.Group("/identity")
	identity.Post("/", identityHandler.AddIdentity)
	identity.Get("/", identityHandler.GetIdentities)

}
