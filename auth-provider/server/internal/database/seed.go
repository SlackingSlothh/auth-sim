package database

import (
	"os"

	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/config"
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/models"
)

func Seed() {
	db := config.GetDB()

	adminName := os.Getenv("ADMIN_NAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPass := os.Getenv("ADMIN_PASS")
	adminGroupName := os.Getenv("ADMIN_GROUP")

	if adminName == "" || adminEmail == "" || adminPass == "" || adminGroupName == "" {
		return
	}

	passwordHash, err := config.HashPassword(adminPass)
	if err != nil {
		return
	}

	adminGroup := models.Group{Name: adminGroupName}
	if err := db.Where("name = ?", adminGroupName).FirstOrCreate(&adminGroup).Error; err != nil {
		return
	}

	admin := models.User{
		Name:         adminName,
		Email:        adminEmail,
		PasswordHash: passwordHash,
		Status:       "active",
	}
	if err := db.Where("email = ?", adminEmail).FirstOrCreate(&admin, admin).Error; err != nil {
		return
	}

	userGroup := models.UserGroup{
		UserID:  admin.ID,
		GroupID: adminGroup.ID,
	}
	if err := db.Where("user_id = ? AND group_id = ?", admin.ID, adminGroup.ID).FirstOrCreate(&userGroup).Error; err != nil {
		return
	}
}