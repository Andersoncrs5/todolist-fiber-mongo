package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"syscall"
	"todolist/config"
	"todolist/handlers"
	"todolist/repositories"
	"todolist/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	config.ConnectDB()
	defer config.CloseDB()

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	app.Use(limiter.New(limiter.Config{
		Max:          100,
		Expiration:   25 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
	}))

	todoRepo := repositories.NewTodoRepository(config.MongoDatabase)
	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	api := app.Group("/api/v1/todos")

	api.Post("/", todoHandler.CreateTodo)
	api.Get("/", todoHandler.GetAllTodos)
	api.Get("/:id", todoHandler.GetTodoByID)
	api.Put("/:id", todoHandler.UpdateTodo)
	api.Delete("/:id", todoHandler.DeleteTodo)
	api.Patch("/:id/toggle", todoHandler.ToggleTodoComplete)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"message": "API is running!",
		})
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		_ = <-c
		log.Println("Shutting down server...")
		_ = app.Shutdown()
	}()

	port := "3000"

	log.Printf("Server is starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Println("Server gracefully stopped.")
}
