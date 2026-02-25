package seed

import (
	"LevelUp_Hub_Backend/internal/modules/rbac"
	"fmt"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {

	defaultRoles := []string{
		"admin",
		"mentor",
		"student",
	}

	for _, roleName := range defaultRoles {

		role:=rbac.Role{Name: roleName}

		err:=db.Where(rbac.Role{Name: roleName}).
		     FirstOrCreate(&role).Error

		if err!=nil{
			fmt.Printf("Error seeding role %s: %v\n", roleName, err)
		}else{
			fmt.Printf("Role ensured: %s\n", roleName)
		}
	}
}
