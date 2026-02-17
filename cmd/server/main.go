package main

import (
	"LevelUp_Hub_Backend/internal/config"
	"LevelUp_Hub_Backend/internal/modules/booking"
	"LevelUp_Hub_Backend/internal/modules/courses"
	"LevelUp_Hub_Backend/internal/modules/message"
	"LevelUp_Hub_Backend/internal/modules/payment"
	"LevelUp_Hub_Backend/internal/modules/profile"
	"LevelUp_Hub_Backend/internal/modules/ratings"
	"LevelUp_Hub_Backend/internal/modules/slot"
	"LevelUp_Hub_Backend/internal/platform/postgres"
	"LevelUp_Hub_Backend/internal/platform/redis"
	"LevelUp_Hub_Backend/internal/routes"
	"log"

	"github.com/gofiber/fiber/v2/middleware/cors"

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
	if err:=db.AutoMigrate(
		&profile.User{},
		&profile.MentorProfile{},
		&courses.Course{},
		&courses.MentorCourse{},
		&slot.MentorSlot{},
		&booking.Booking{},
		&message.Conversation{},
		&message.Message{},
		&payment.Payment{},
		&payment.Wallet{},
		&payment.WalletTransaction{},
		&payment.WithdrawRequest{},
		&ratings.Rating{},
	);err!=nil{
		log.Fatal(err)
	}

	// Connect Redis
	rdb, err := redis.NewRedisClient(cfg)
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}
	_ = rdb

	// Create Fiber app
	app := fiber.New()

	//for frontend connect
	app.Use(cors.New(cors.Config{
    AllowOrigins: "http://localhost:5173",
    AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
    AllowHeaders: "Origin, Content-Type, Accept,Authorization",
    AllowCredentials: true,
}))

	//setup routes
	routes.SetUp(app,db,rdb,cfg.JWTSecret,cfg.RazorpayClient,cfg.RazorpayKey)

	// Health check route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(" LevelUp Hub backend running")
	})

	log.Println("Server running on port", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}