package routes

import (
	"github.com/gorilla/mux"
	"github.com/jeffemart/Gotham/handlers"
	"github.com/jeffemart/Gotham/middlewares"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Rotas públicas
	router.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}", handlers.GetUser).Methods("GET")
	router.HandleFunc("/users", handlers.CreateUser).Methods("POST")

	// Rota de login (pública)
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	
	// Rotas protegidas (somente para administradores)
	adminRoutes := router.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middlewares.RoleMiddleware("admin"))
	// adminRoutes.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	adminRoutes.HandleFunc("/users/{id:[0-9]+}", handlers.UpdateUser).Methods("PUT")
	adminRoutes.HandleFunc("/users/{id:[0-9]+}", handlers.DeleteUser).Methods("DELETE")

	// Rotas para qualquer usuário autenticado (como exemplo, apenas admin ou agente podem acessar)
	protectedRoutes := router.PathPrefix("/protected").Subrouter()
	protectedRoutes.Use(middlewares.RoleMiddleware("admin", "agente"))
	protectedRoutes.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")

	return router
}
