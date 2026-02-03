package main

import (
	"LevelUp_Hub_Backend/internal/config"
	"LevelUp_Hub_Backend/internal/modules/user"
	"LevelUp_Hub_Backend/internal/platform/postgres"
	"LevelUp_Hub_Backend/internal/platform/redis"
	"LevelUp_Hub_Backend/internal/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	//Load config
	cfg := config.LeadConfig()

	//Connect postgres
	db,err:=postgres.NewPostgresConnection(cfg)
	if err !=nil{
		log.Fatal("Postgres connection failed:",err)
	}
	//run migrations
	db.AutoMigrate(
		&user.User{},
	)

	// Connect Redis
	rdb, err := redis.NewRedisClient(cfg)
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}
	_ = rdb

	// Create Fiber app
	app := fiber.New()

	//setup routes
	routes.SetUp(app,db)

	// Health check route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(" LevelUp Hub backend running")
	})

	log.Println("Server running on port", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}