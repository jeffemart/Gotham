package routes

import (
	"github.com/gorilla/mux"
	"github.com/jeffemart/Gotham/internal/handlers"
	"github.com/jeffemart/Gotham/internal/middlewares"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"github.com/jeffemart/Gotham/internal/models"
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

	// Rotas protegidas com capacidades específicas
	adminRoutes := r.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middlewares.AuthMiddleware)
	
	// Rotas de usuário com capacidades específicas
	adminRoutes.Handle("/users/{id:[0-9]+}", 
		middlewares.CapabilityMiddleware(models.CapabilityUpdateUser)(
			http.HandlerFunc(handlers.UpdateUser),
		)).Methods("PUT")
	
	adminRoutes.Handle("/users/{id:[0-9]+}", 
		middlewares.CapabilityMiddleware(models.CapabilityDeleteUser)(
			http.HandlerFunc(handlers.DeleteUser),
		)).Methods("DELETE")

	// Rotas de tarefas
	protectedRoutes := r.PathPrefix("/protected").Subrouter()
	protectedRoutes.Use(middlewares.AuthMiddleware)
	
	protectedRoutes.Handle("/tasks", 
		middlewares.CapabilityMiddleware(models.CapabilityViewTasks)(
			http.HandlerFunc(handlers.GetTasks),
		)).Methods("GET")
}
