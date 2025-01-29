package routes

import (
	"github.com/gorilla/mux"
	"github.com/jeffemart/Gotham/internal/handlers"
	"github.com/jeffemart/Gotham/internal/middlewares"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(r *mux.Router) {
	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8080",
			"http://localhost:8000",
			"http://localhost",
			"http://127.0.0.1",
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowedHeaders: []string{
			"Content-Type", "Authorization", "X-Requested-With",
		},
		AllowCredentials: true,
	})
	r.Use(c.Handler)

	// Endpoint da documentação Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Rotas públicas
	// @Summary Lista todos os usuários
	// @Description Retorna todos os usuários registrados no sistema
	// @Tags users
	// @Produce json
	// @Success 200 {array} map[string]interface{}
	// @Router /users [get]
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")

	// @Summary Retorna um usuário pelo ID
	// @Description Busca um usuário pelo ID fornecido
	// @Tags users
	// @Param id path int true "ID do usuário"
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /users/{id} [get]
	r.HandleFunc("/users/{id:[0-9]+}", handlers.GetUser).Methods("GET")

	// @Summary Cria um novo usuário
	// @Description Adiciona um usuário ao sistema
	// @Tags users
	// @Accept json
	// @Produce json
	// @Param user body map[string]interface{} true "Dados do usuário"
	// @Success 201 {object} map[string]interface{}
	// @Router /users [post]
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")

	// @Summary Realiza login
	// @Description Autentica um usuário e retorna o token JWT
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param credentials body map[string]interface{} true "Credenciais de login"
	// @Success 200 {object} map[string]string
	// @Router /login [post]
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	// Rotas protegidas (somente para administradores)
	adminRoutes := r.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middlewares.AuthMiddleware)
	adminRoutes.Use(middlewares.RoleMiddleware("admin"))

	// @Summary Atualiza um usuário
	// @Description Atualiza os dados de um usuário pelo ID
	// @Tags admin
	// @Param id path int true "ID do usuário"
	// @Accept json
	// @Produce json
	// @Param user body map[string]interface{} true "Dados atualizados do usuário"
	// @Success 200 {object} map[string]interface{}
	// @Router /admin/users/{id} [put]
	adminRoutes.HandleFunc("/users/{id:[0-9]+}", handlers.UpdateUser).Methods("PUT")

	// @Summary Remove um usuário
	// @Description Deleta um usuário pelo ID
	// @Tags admin
	// @Param id path int true "ID do usuário"
	// @Produce json
	// @Success 204 "No Content"
	// @Router /admin/users/{id} [delete]
	adminRoutes.HandleFunc("/users/{id:[0-9]+}", handlers.DeleteUser).Methods("DELETE")

	// Rotas protegidas para usuários autenticados (admin ou agente)
	protectedRoutes := r.PathPrefix("/protected").Subrouter()
	protectedRoutes.Use(middlewares.AuthMiddleware)
	protectedRoutes.Use(middlewares.RoleMiddleware("admin", "agente"))

	// @Summary Lista todas as tarefas
	// @Description Retorna as tarefas disponíveis para o usuário autenticado
	// @Tags protected
	// @Produce json
	// @Success 200 {array} map[string]interface{}
	// @Router /protected/tasks [get]
	protectedRoutes.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
}
