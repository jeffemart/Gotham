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

	// Criar roles com capacidades
	adminRole := models.Role{
		Name:        "admin",
		Permissions: permissions,
		Capabilities: []string{
			"*", // Admin tem todas as capacidades
		},
	}
	database.DB.Create(&adminRole)

	agentRole := models.Role{
		Name:        "agent",
		Permissions: permissions,
		Capabilities: []string{
			models.CapabilityReadUser,
			models.CapabilityUpdateUser,
			models.CapabilityViewTasks,
			models.CapabilityManageTasks,
		},
	}
	database.DB.Create(&agentRole)

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
