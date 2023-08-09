package user

import "github.com/gofiber/fiber/v2"

type UserController struct {
	storage *UserStorage
}

func NewUserController(storage *UserStorage) *UserController {
	return &UserController{
		storage: storage,
	}
}

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type createUserResponse struct {
	Message string `json:"message"`
}

func (u *UserController) profile(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "You are logged in",
	})
}

// @Summary Create user.
// @Description creates user.
// @Tags users
// @Accept */*
// @Produce json
// @Param user body createUserRequest true "User to create"
// @Success 200 {object} createUserResponse
// @Router /users [post]
func (t *UserController) create(c *fiber.Ctx) error {
	// parse the request body
	var req createUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// TODO: check if user does not exist or value changed

	// create the user
	message, err := t.storage.createUser(req.Name, req.Email, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(
		createUserResponse{
			message,
		})
}
