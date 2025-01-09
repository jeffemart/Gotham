package database

import (
	"fmt"
	"log"

	"github.com/jeffemart/Gotham/settings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect estabelece a conexão com o banco de dados
func Connect() {
	// Carregar as configurações das variáveis de ambiente
	config := settings.LoadSettings()

	// Criar o Data Source Name (DSN) para a conexão com o banco
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.Host,
		config.Database.Port,
	)

	// Conectar ao banco de dados
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	log.Println("Conexão com o banco de dados estabelecida.")
}
