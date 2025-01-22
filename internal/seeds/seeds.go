package seeds

import (
	"github.com/jeffemart/Gotham/internal/database"
	"github.com/jeffemart/Gotham/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// SeedDatabase popula o banco de dados com dados iniciais
func SeedDatabase() error {
	// Criar permissões
	permissions := []models.Permission{
		{Name: "view_tasks"},
		{Name: "edit_tasks"},
		{Name: "delete_tasks"},
	}

	for _, permission := range permissions {
		database.DB.Create(&permission)
	}

	// Criar roles
	adminRole := models.Role{
		Name:        "admin",
		Permissions: permissions,
	}
	database.DB.Create(&adminRole)

	// Criar usuário admin
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	adminUser := models.User{
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: string(hashedPassword),
		RoleID:   adminRole.ID,
	}
	database.DB.Create(&adminUser)

	return nil
}
