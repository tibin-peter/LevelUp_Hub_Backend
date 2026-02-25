package seed

import (
	"LevelUp_Hub_Backend/internal/modules/rbac"

	"gorm.io/gorm"
)

var defaultPermissions = []rbac.Permission{
  
	//student

	{Slug: "dashboard", Name: "Dashboard"},
	{Slug: "connected_mentors", Name: "Connected Mentors"},
	{Slug: "courses", Name: "Courses"},
	{Slug: "bookings", Name: "Bookings"},
	{Slug: "messages", Name: "Messages"},
	{Slug: "payments", Name: "Payments"},
	{Slug: "settings", Name: "Settings"},

	// Mentor
	{Slug: "dashboard", Name: "Dashboard"},
	{Slug: "explore_courses", Name: "Explore Courses"},
	{Slug: "my_courses", Name: "My Courses"},
	{Slug: "sessions", Name: "Sessions"},
	{Slug: "earnings", Name: "Earnings"},
	{Slug: "booking_requests", Name: "Booking Requests"},

	// Admin
	{Slug: "students", Name: "Students"},
	{Slug: "mentors", Name: "Mentors"},
	{Slug: "mentor_approvals", Name: "Mentor Approvals"},
	{Slug: "complaints", Name: "Complaints"},
	{Slug: "wallet", Name: "Wallet"},
	{Slug: "profile", Name: "Profile"},
}

func SeedPermissions(db *gorm.DB) error {

	for _, p := range defaultPermissions {

		if err := db.
			Where("slug = ?", p.Slug).
			FirstOrCreate(&rbac.Permission{}, p).Error; err != nil {
			return err
		}
	}

	return nil
}

func SeedAdminPermissions(db *gorm.DB) error {

	var admin rbac.Role
	if err := db.
		Where("name = ?", "admin").
		First(&admin).Error; err != nil {
		return err
	}

	var perms []rbac.Permission
	if err := db.Find(&perms).Error; err != nil {
	return err
}

	// attach all permissions
	return db.
		Model(&admin).
		Association("Permissions").
		Append(&perms)
		// Replace(&perms)
}