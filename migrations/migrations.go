package migrations

import (
	"log"

	"github.com/jeffemart/Gotham/database"
	"github.com/jeffemart/Gotham/models"
)

// Run executa a migração do banco de dados
func Run() {
	// Migrar a tabela User
	err := database.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Erro ao migrar a tabela: %v", err)
	}

	log.Println("Tabela 'users' criada com sucesso!")
}
