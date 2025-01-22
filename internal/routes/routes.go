package routes

import (
	"github.com/gorilla/mux"
	"github.com/jeffemart/Gotham/internal/handlers"
	"github.com/jeffemart/Gotham/internal/middlewares"
	"github.com/rs/cors"
)

func SetupRoutes() *mux.Router {
	// Configurar CORS com a biblioteca
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8080", // Permitindo o acesso do localhost:8080
			"http://localhost:8000", // Permitindo o acesso do localhost:8000
			"http://localhost",      // Permitindo o acesso do localhost
			"http://127.0.0.1",      // Permitindo o acesso de 127.0.0.1
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowedHeaders: []string{
			"Content-Type", "Authorization", "X-Requested-With",
		},
		AllowCredentials: true,
	})

	router := mux.NewRouter()

	// Aplicar o middleware CORS
	router.Use(c.Handler)

	// Rotas públicas
	router.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}", handlers.GetUser).Methods("GET")
	router.HandleFunc("/users", handlers.CreateUser).Methods("POST")

	// Rota de login (pública)
	router.HandleFunc("/login", handlers.Login).Methods("POST")

	// Rotas protegidas (somente para administradores)
	adminRoutes := router.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middlewares.AuthMiddleware)          // Primeiro verifica autenticação
	adminRoutes.Use(middlewares.RoleMiddleware("admin")) // Depois verifica a role
	adminRoutes.HandleFunc("/users/{id:[0-9]+}", handlers.UpdateUser).Methods("PUT")
	adminRoutes.HandleFunc("/users/{id:[0-9]+}", handlers.DeleteUser).Methods("DELETE")

	// Rotas para qualquer usuário autenticado (exemplo de admin ou agente)
	protectedRoutes := router.PathPrefix("/protected").Subrouter()
	protectedRoutes.Use(middlewares.AuthMiddleware)                    // Primeiro verifica autenticação
	protectedRoutes.Use(middlewares.RoleMiddleware("admin", "agente")) // Depois verifica a role
	protectedRoutes.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")

	return router
}
