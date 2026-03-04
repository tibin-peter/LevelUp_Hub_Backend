package seed

import (
	"LevelUp_Hub_Backend/internal/modules/profile"
	"LevelUp_Hub_Backend/internal/utils"
	"fmt"
	"gorm.io/gorm"
)

func Addadmin(db *gorm.DB) {

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

  db.FirstOrCreate(&admin)

  fmt.Println("Admin seeded successfully")
}