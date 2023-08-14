package order

import (
	"github.com/gofiber/fiber/v2"
)

type OrderController struct {
	storage *OrderStorage
}

func NewOrderController(storage *OrderStorage) *OrderController {
	return &OrderController{
		storage: storage,
	}
}

type createOrderRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type createOrderResponse struct {
	Message string `json:"message"`
}

// @Summary Create one todo.
// @Description creates one todo.
// @Tags todos
// @Accept */*
// @Produce json
// @Param todo body createOrderRequest true "Order to create"
// @Success 200 {object} createOrderResponse
// @Router /todos [post]
func (t *OrderController) create(c *fiber.Ctx) error {
	// parse the request body
	var req createOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// create the todo
	message, err := t.storage.createOrder(req.Title, req.Description, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create todo",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(
		createOrderResponse{
			message,
		})
}

// @Summary Get all todos.
// @Description fetch every todo available.
// @Tags todos
// @Accept */*
// @Produce json
// @Success 200 {object} []todoDB
// @Router /todos [get]
func (t *OrderController) getAll(c *fiber.Ctx) error {
	// get all todos
	todos, err := t.storage.getAllOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get todos",
		})
	}

	return c.JSON(todos)
}
