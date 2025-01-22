package migrations

import (
	"log"

	"github.com/jeffemart/Gotham/internal/database"
	"github.com/jeffemart/Gotham/internal/models"
)

// Run executa a migração do banco de dados
// Run executa as migrações do banco de dados
func Run() error {
	// Obter a instância do banco de dados
	db := database.DB

	log.Println("Iniciando migrações...")

	// Executar a migração da tabela `users`
	if err := db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.RolePermission{}); err != nil {
		log.Printf("Erro ao executar migração da tabela `users`: %v\n", err)
		return err
	}

	log.Println("Migrações concluídas com sucesso!")
	return nil
}
