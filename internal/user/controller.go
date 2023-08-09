package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

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

type auth0UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type getUserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *UserController) profile(c *fiber.Ctx) error {
	authToken := c.Locals("authTokenPart").(string)
	auth0Url := c.Locals("auth0Url").(string)
	req, _ := http.NewRequest("GET", auth0Url, nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get user profile",
		})
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
	var response auth0UserResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse user profile",
		})
	}

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
