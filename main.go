package main

import (
	"fmt"
	"log"

	"github.com/jeffemart/Gotham/database"
	"github.com/jeffemart/Gotham/migrations"
	"github.com/jeffemart/Gotham/settings"
)

func main() {
	// Carregar configurações
	config := settings.LoadSettings()

	// Exemplo de uso
	fmt.Println("Conectando ao banco de dados:")
	fmt.Printf("Driver: %s\nHost: %s\nPort: %s\nUser: %s\n",
		config.Database.Driver,
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
	)

	// Conectar ao banco de dados
	database.Connect()

	// Executar a migração
	migrations.Run()

	log.Println("Migração concluída com sucesso!")
}
