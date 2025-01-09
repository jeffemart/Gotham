package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jeffemart/Gotham/database"
	"github.com/jeffemart/Gotham/migrations"
	"github.com/jeffemart/Gotham/routes"
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
	if err := migrations.Run(); err != nil {
		log.Fatalf("Erro ao executar migração: %v", err)
	}

	log.Println("Migração concluída e conexão fechada com sucesso!")

	// Configurar rotas
	router := routes.SetupRoutes()

	log.Println("Servidor iniciado na porta 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
