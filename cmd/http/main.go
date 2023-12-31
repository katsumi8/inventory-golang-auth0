package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/katsumi/inventory_api/config"
	_ "github.com/katsumi/inventory_api/docs"
	"github.com/katsumi/inventory_api/internal/order"
	"github.com/katsumi/inventory_api/internal/storage"
	"github.com/katsumi/inventory_api/internal/todo"
	"github.com/katsumi/inventory_api/internal/user"
	"github.com/katsumi/inventory_api/pkg/shutdown"
)

// @title Tapir App Template
// @version 2.0
// @description An example template of a Golang backend API using Fiber and MongoDB
// @contact.name Ben Davis
// @license.name MIT
// @BasePath /
func main() {
	// setup exit code for graceful shutdown
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// load config
	env, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	// run the server
	cleanup, err := run(env)

	// run the cleanup after the server is terminated
	defer cleanup()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	// ensure the server is shutdown gracefully & app runs
	shutdown.Gracefully()
}

func run(env config.EnvVars) (func(), error) {
	app, cleanup, err := buildServer(env)
	if err != nil {
		return nil, err
	}

	// start the server
	go func() {
		app.Listen("0.0.0.0:" + env.PORT)
	}()

	// return a function to close the server and database
	return func() {
		cleanup()
		app.Shutdown()
	}, nil
}

func buildServer(env config.EnvVars) (*fiber.App, func(), error) {
	postgreConfig := storage.Config{
		Host:     env.POSTGRES_HOST,
		Port:     env.POSTGRES_PORT,
		Password: env.POSTGRES_PASSWORD,
		User:     env.POSTGRES_USER,
		DBName:   env.POSTGRES_NAME,
		SSLMode:  env.POSTGRES_SSL, // SSLModeもEnvVarsに追加する必要があります
	}

	// init the storage
	db, err := storage.NewConnection(postgreConfig)
	if err != nil {
		return nil, nil, err
	}

	err = storage.Migrate(db)
	if err != nil {
		log.Fatalf("Could not migrate users: %v", err)
	}

	// create the fiber app
	app := fiber.New()

	// add middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     env.FRONTEND_ORIGIN, // これはフロントエンドのアドレスです
		AllowCredentials: true,
	}))
	app.Use(logger.New())

	// add health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	// add docs
	app.Get("/swagger/*", swagger.HandlerDefault)

	// create the user domain
	userStore := user.NewUserStorage(db)
	userController := user.NewUserController(userStore)
	user.CreateUserGroup(app, userController, env)

	// create the todo domain
	todoStore := todo.NewTodoStorage(db)
	todoController := todo.NewTodoController(todoStore)
	todo.AddTodoRoutes(app, todoController)

	// create the order domain
	orderStore := order.NewOrderStorage(db)
	orderController := order.NewOrderController(orderStore)
	order.AddOrderRoutes(app, orderController)

	return app, func() {
		storage.CloseConnection(db)
	}, nil
}
