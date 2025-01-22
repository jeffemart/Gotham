package routes

import (
	"github.com/gorilla/mux"
	"github.com/jeffemart/Gotham/internal/handlers"
	"github.com/jeffemart/Gotham/internal/middlewares"
	"github.com/rs/cors"
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

	// Aplicar o middleware CORS
	r.Use(c.Handler)

	// Rotas públicas
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	// Rotas protegidas (somente para administradores)
	adminRoutes := r.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middlewares.AuthMiddleware)
	adminRoutes.Use(middlewares.RoleMiddleware("admin"))
	adminRoutes.HandleFunc("/users/{id:[0-9]+}", handlers.UpdateUser).Methods("PUT")
	adminRoutes.HandleFunc("/users/{id:[0-9]+}", handlers.DeleteUser).Methods("DELETE")

	// Rotas para qualquer usuário autenticado (admin ou agente)
	protectedRoutes := r.PathPrefix("/protected").Subrouter()
	protectedRoutes.Use(middlewares.AuthMiddleware)
	protectedRoutes.Use(middlewares.RoleMiddleware("admin", "agente"))
	protectedRoutes.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
}
