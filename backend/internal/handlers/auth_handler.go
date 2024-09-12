package handlers

import (
	"github.com/ashmitsharp/msg-agg/internal/services"
	"github.com/ashmitsharp/msg-agg/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	UserService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{UserService: userService}
}

func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.ValidateEmail(input.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := utils.ValidatePassword(input.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := utils.ValidateName(input.Name); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	matrixService := services.NewMatrixService("http://synapse:8008")

	matrixUser := utils.GenerateMatrixUsername(input.Email)
	matrixPassword, err := utils.GenerateSecurePassword()
	if err != nil {
		return err
	}

	matrixResp, err := matrixService.CreateAccount(matrixUser, matrixPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Matrix account creation failed",
		})
	}

	user, err := h.UserService.CreateUser(input.Email, input.Password, input.Name, matrixPassword, matrixResp)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": token})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := h.UserService.AuthenticateUser(input.Email, input.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	slackIdentity, err := h.UserService.GetIdentifier(user.MasterUserID, "slack")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Platforms Credential Not Found"})
	}

	matrixService := services.NewMatrixService("http://synapse:8008")

	client, err := matrixService.Login(user.MatrixUserID, user.MatrixAccessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Matrix login failed",
		})
	}

	// Initialize the BridgeService and connect to linked bridges
	bridgeService := services.NewBridgeService(client)
	err = bridgeService.ConnectSlackBridge(slackIdentity.Identifier)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to connect to Slack bridge",
		})
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    token, // Store the JWT token in the cookie
		Path:     "/",
		HTTPOnly: true, // Ensure cookie is not accessible via JavaScript
		Secure:   true, // Ensure the cookie is only sent over HTTPS (disable in dev if needed)
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":          token,
		"matrix_user_id": user.MatrixUserID,
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	newToken, err := utils.GenerateToken(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate new token"})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    newToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
	})
	return c.JSON(fiber.Map{"token": newToken})
}
