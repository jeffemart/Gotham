package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeffemart/Gotham/internal/database"
	"github.com/jeffemart/Gotham/internal/routes"
	"github.com/jeffemart/Gotham/internal/settings"
	"github.com/jeffemart/Gotham/migrations"
)

//go:generate swag init -g cmd/gotham/main.go -o docs

func main() {
	// @title           Gotham API
	// @version         1.0.0
	// @description     API para gerenciamento de usuários e autenticação
	// @termsOfService  http://swagger.io/terms/

	// @contact.name   Jefferson Martins
	// @contact.url    https://www.linkedin.com/in/jefferson-martins-a6802b249/
	// @contact.email  jefferson.developers@gmail.com

	// @license.name  MIT
	// @license.url   https://opensource.org/licenses/MIT

	// @host      localhost:8000
	// @BasePath  /api/v1

	// @securityDefinitions.apikey Bearer
	// @in header
	// @name Authorization
	// @description Type "Bearer" followed by a space and JWT token.

	// @tag.name users
	// @tag.description Operações relacionadas a usuários

	// @tag.name auth
	// @tag.description Operações de autenticação

	// Carregar configurações
	config := settings.LoadSettings()

	// Conectar ao banco de dados
	database.Connect()

	// Executar a migração
	if err := migrations.Run(); err != nil {
		log.Fatalf("Erro ao executar migração: %v", err)
	}

	// Criar o roteador principal
	r := mux.NewRouter()

	// Setup das rotas
	routes.SetupRoutes(r)

	// Iniciar servidor na porta configurada
	port := config.App.Port
	log.Printf("Servidor iniciado em http://localhost:%s", port)
	log.Printf("Documentação Swagger disponível em http://localhost:%s/swagger/", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
