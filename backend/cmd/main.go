package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"task-manager/config"
	"task-manager/database"
	"task-manager/queue"
	"task-manager/routes"
)

func main() {
	_ = godotenv.Load()
	database.ConnectDB()
	queue.InitAsynq()
	queue.InitAsynqClient(config.RedisUrl)

	app := fiber.New()
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
