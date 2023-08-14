package user

import (
	"encoding/json"
	"io"
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
	Sub           string `json:"sub"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Nickname      string `json:"nickname"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	UpdatedAt     string `json:"updated_at"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type getUserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *UserController) getAuth0User(c *fiber.Ctx) (*auth0UserResponse, error) {
	authToken, ok := c.Locals("validatedAccessToken").(string)
	if !ok || authToken == "" {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "authToken is missing or incorrect type",
		})
	}

	auth0Url, ok := c.Locals("auth0Url").(string)
	if !ok || auth0Url == "" {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "auth0Url is missing or incorrect type",
		})
	}

	req, _ := http.NewRequest("GET", auth0Url, nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get user profile",
		})
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var user auth0UserResponse
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse user profile",
		})
	}

	return &user, nil
}

func (u *UserController) profile(c *fiber.Ctx) error {
	_, err := u.getAuth0User(c)
	if err != nil {
		return err
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
func (u *UserController) create(c *fiber.Ctx) error {
	user, err := u.getAuth0User(c)
	if err != nil {
		return err
	}

	// check if user does not exist or value changed
	doesUserExist, err := u.storage.IsUserWithEmailPresent(user.Email)
	if doesUserExist {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already exists",
		})
	}

	// create the user
	_, err = u.storage.createUser(user.Name, user.Email, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(
		createUserResponse{
			Message: "User created",
		})
}

// @Summary Get user.
// @Description gets user.
// @Tags users
// @Accept */*
// @Produce json
// @Param email path string true "Email of user to get"
// @Success 200 {object} getUserResponse
// @Router /users/{email} [get]
func (u *UserController) getMe(c *fiber.Ctx) error {
	user, err := u.getAuth0User(c)
	if err != nil {
		return err
	}

	userFromDb, err := u.storage.getUserByEmail(user.Email)
	if err != nil {
		if err.Error() == "USER_NOT_FOUND" {
			// create the user
			createdUser, err := u.storage.createUser(user.Name, user.Email, false)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to get user",
				})
			}
			return c.JSON(
				getUserResponse{
					Name:  createdUser.Username,
					Email: createdUser.Email,
				})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get user",
		})
	}

	return c.JSON(
		getUserResponse{
			Name:  userFromDb.Username,
			Email: userFromDb.Email,
		})
}
