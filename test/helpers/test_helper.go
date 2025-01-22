package helpers

import (
	"github.com/jeffemart/Gotham/internal/database"
	"github.com/jeffemart/Gotham/internal/models"
)

func SetupTestDB() {
	database.Connect()
	CleanupTestDB()
}

func CleanupTestDB() {
	database.DB.Exec("DELETE FROM users")
	database.DB.Exec("DELETE FROM roles")
	database.DB.Exec("DELETE FROM permissions")
	database.DB.Exec("DELETE FROM role_permissions")
}

func CreateTestUser() models.User {
	user := models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "Test123!",
		RoleID:   1,
	}
	database.DB.Create(&user)
	return user
}
