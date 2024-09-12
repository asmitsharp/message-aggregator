package handlers

import (
	"github.com/ashmitsharp/msg-agg/internal/services"
	"github.com/gofiber/fiber/v2"
)

type IdentityHandler struct {
	UserService *services.UserService
}

func NewIdentityHandler(userService *services.UserService) *IdentityHandler {
	return &IdentityHandler{UserService: userService}
}

func (h *IdentityHandler) AddIdentity(c *fiber.Ctx) error {
	var input struct {
		Platform   string `json:"platform"`
		Identifier string `json:"identifier"`
		Credential string `json:"credential"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	userID := c.Locals("userID").(string)

	identity, err := h.UserService.AddIdentity(userID, input.Platform, input.Identifier, input.Credential)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add identity"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"identity": identity})
}

func (h *IdentityHandler) GetIdentities(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	identities, err := h.UserService.GetIdentities(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve identities"})
	}

	return c.JSON(fiber.Map{"identities": identities})
}
