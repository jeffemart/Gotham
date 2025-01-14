package main

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/jeffemart/Gotham/database"
	"github.com/jeffemart/Gotham/migrations"
	"github.com/jeffemart/Gotham/models"
	"github.com/jeffemart/Gotham/routes"
)

func main() {
	// @title Gotham API
	// @version 1.1.10
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

	log.Println("Servidor iniciado na porta 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
