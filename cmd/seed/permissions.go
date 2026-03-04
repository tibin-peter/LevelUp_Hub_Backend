package seed

import (
	"LevelUp_Hub_Backend/internal/modules/rbac"
	"gorm.io/gorm"
)

var defaultPermissions = []rbac.Permission{
	// Student
	{Slug: "dashboard", Name: "Dashboard"},
	{Slug: "explore_mentors", Name: "Explore Mentors"},
	{Slug: "mentors", Name: "Connected Mentors"},
	{Slug: "courses", Name: "Courses"},
	{Slug: "bookings", Name: "Bookings"},
	{Slug: "messages", Name: "Messages"},
	{Slug: "payments", Name: "Payments"},
	{Slug: "settings", Name: "Settings"},

	// Mentor
	{Slug: "explore_courses", Name: "Explore Courses"},
	{Slug: "my_courses", Name: "My Courses"},
	{Slug: "sessions", Name: "Sessions"},
	{Slug: "earnings", Name: "Earnings"},
	{Slug: "booking_requests", Name: "Booking Requests"},

	// Admin
	{Slug: "students", Name: "Students"},
	{Slug: "admin_mentors", Name: "Mentors"},
	{Slug: "mentor_approvals", Name: "Mentor Approvals"},
	{Slug: "complaints", Name: "Complaints"},
	{Slug: "wallet", Name: "Wallet"},
	{Slug: "profile", Name: "Profile"},
}

func SeedPermissions(db *gorm.DB) error {
	for _, p := range defaultPermissions {
		if err := db.Where("slug = ?", p.Slug).FirstOrCreate(&rbac.Permission{}, p).Error; err != nil {
			return err
		}
	}
	return nil
}

func SeedRolePermissions(db *gorm.DB) error {
	// 1. Admin gets all
	var admin rbac.Role
	if err := db.Where("name = ?", "admin").First(&admin).Error; err == nil {
		var all []rbac.Permission
		db.Find(&all)
		db.Model(&admin).Association("Permissions").Replace(&all)
	}

	// 2. Student Permissions
	var student rbac.Role
	if err := db.Where("name = ?", "student").First(&student).Error; err == nil {
		slugs := []string{"dashboard", "explore_mentors", "mentors", "courses", "bookings", "messages", "payments", "settings"}
		var perms []rbac.Permission
		db.Where("slug IN ?", slugs).Find(&perms)
		db.Model(&student).Association("Permissions").Replace(&perms)
	}

	// 3. Mentor Permissions
	var mentor rbac.Role
	if err := db.Where("name = ?", "mentor").First(&mentor).Error; err == nil {
		slugs := []string{"dashboard", "explore_courses", "my_courses", "sessions", "earnings", "messages", "booking_requests", "settings"}
		var perms []rbac.Permission
		db.Where("slug IN ?", slugs).Find(&perms)
		db.Model(&mentor).Association("Permissions").Replace(&perms)
	}

	return nil
}