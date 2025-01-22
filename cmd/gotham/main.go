package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jeffemart/Gotham/internal/database"
	"github.com/jeffemart/Gotham/internal/models"
	"github.com/jeffemart/Gotham/internal/routes"
	"github.com/jeffemart/Gotham/internal/settings"
	"github.com/jeffemart/Gotham/migrations"
	"golang.org/x/crypto/bcrypt"
)

func serveSwaggerFile(w http.ResponseWriter, r *http.Request) {
	// Lê o arquivo JSON do Swagger
	jsonFile, err := os.ReadFile("docs/openapi.json")
	if err != nil {
		log.Printf("Erro ao ler arquivo Swagger: %v", err)
		http.Error(w, "Não foi possível ler o arquivo Swagger", http.StatusInternalServerError)
		return
	}

	// Define o tipo de conteúdo como JSON
	w.Header().Set("Content-Type", "application/json")
	// Permite CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Escreve o conteúdo do arquivo
	w.Write(jsonFile)
}

func main() {
	// @title Gotham API
	// @version 1.1.12
	// @description Gotham é um projeto de uma API desenvolvido para gerenciar usuários, permissões e autenticação de forma robusta e segura.

	// @contact.name Jefferson Martins
	// @contact.url https://www.linkedin.com/in/jefferson-martins-a6802b249/
	// @contact.email jefferson.developers@gmail.com

	// @license.name MIT
	// @license.url https://opensource.org/licenses/MIT

	// @host localhost:8000
	// @BasePath /

	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization

	// Carregar configurações
	config := settings.LoadSettings()

	// Conectar ao banco de dados
	database.Connect()

	// Executar a migração
	if err := migrations.Run(); err != nil {
		log.Fatalf("Erro ao executar migração: %v", err)
	}

	log.Println("Migração concluída e conexão fechada com sucesso!")

	// Criar permissões
	permission1 := models.Permission{Name: "view_tasks"}
	permission2 := models.Permission{Name: "edit_tasks"}
	database.DB.Create(&permission1)
	database.DB.Create(&permission2)

	// Criar roles e associar permissões
	adminRole := models.Role{Name: "admin", Permissions: []models.Permission{permission1, permission2}}
	agentRole := models.Role{Name: "agente", Permissions: []models.Permission{permission1}}
	database.DB.Create(&adminRole)
	database.DB.Create(&agentRole)

	// Função para criptografar a senha
	hashPassword := func(password string) (string, error) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}
		return string(hashedPassword), nil
	}

	// Criar usuários e associar roles com a senha criptografada
	adminPassword, err := hashPassword("admin123")
	if err != nil {
		log.Fatalf("Erro ao criptografar senha do administrador: %v", err)
	}
	agentPassword, err := hashPassword("agent123")
	if err != nil {
		log.Fatalf("Erro ao criptografar senha do agente: %v", err)
	}

	adminUser := models.User{Name: "Admin User", Email: "admin@example.com", Password: adminPassword, RoleID: adminRole.ID}
	agentUser := models.User{Name: "Agent User", Email: "agent@example.com", Password: agentPassword, RoleID: agentRole.ID}
	database.DB.Create(&adminUser)
	database.DB.Create(&agentUser)

	// Configurar rotas
	router := routes.SetupRoutes()

	// Adicionar rota para servir o arquivo Swagger JSON
	router.HandleFunc("/swagger/openapi.json", serveSwaggerFile).Methods("GET")

	// Servir o Swagger UI
	fs := http.FileServer(http.Dir("api/swagger"))
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", fs))

	// Iniciar servidor na porta configurada
	port := config.App.Port
	log.Printf("Servidor iniciado na porta %s...", port)
	log.Printf("Documentação Swagger disponível em http://localhost:%s/swagger/", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
