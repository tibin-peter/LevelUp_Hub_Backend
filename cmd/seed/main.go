package main

import (
	"LevelUp_Hub_Backend/internal/config"
	"LevelUp_Hub_Backend/internal/modules/profile"
	"LevelUp_Hub_Backend/internal/platform/postgres"
	"LevelUp_Hub_Backend/internal/utils"
	"fmt"
	"log"
)

func main() {
  //Load config
	cfg := config.LeadConfig()

	//Connect postgres
	db,err:=postgres.NewPostgresConnection(cfg)
	if err !=nil{
		log.Fatal("Postgres connection failed:",err)
	}

  password := "admin123"

  hashed, _ := utils.HashPassword(password)

  admin := profile.User{
    Name: "Admin",
    Email: "admin@gmail.com",
    Password: hashed,
		IsVerified: true,
		ProfilePicURL: "https://img.freepik.com/free-vector/follow-me-social-business-theme-design_24877-50426.jpg?semt=ais_user_personalization&w=740&q=80",
    Role: "admin",
  }

  db.Create(&admin)

  fmt.Println("Admin seeded successfully")
}