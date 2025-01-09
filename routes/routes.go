package routes

import (
	"net/http"

	"github.com/jeffemart/Gotham/handlers"

	"github.com/gorilla/mux"
)

// SetupRoutes configura as rotas para a aplicação
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Rotas de usuário
	router.HandleFunc("/users", handlers.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", handlers.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", handlers.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", handlers.DeleteUser).Methods(http.MethodDelete)

	return router
}
