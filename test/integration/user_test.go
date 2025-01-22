package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeffemart/Gotham/internal/database"
	"github.com/jeffemart/Gotham/internal/models"
	"github.com/jeffemart/Gotham/internal/routes"
	"github.com/jeffemart/Gotham/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestUserCreation(t *testing.T) {
	// Setup
	helpers.SetupTestDB()
	defer helpers.CleanupTestDB()

	// Inicializar o router
	router := routes.SetupRoutes()

	// Criar payload do usuário
	user := models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "Test123!",
		RoleID:   2,
	}

	payload, _ := json.Marshal(user)

	// Criar request
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	// Criar response recorder
	w := httptest.NewRecorder()

	// Executar request
	router.ServeHTTP(w, req)

	// Verificar status code
	assert.Equal(t, http.StatusCreated, w.Code)

	// Verificar resposta
	var response models.User
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, response.Name)
	assert.Equal(t, user.Email, response.Email)

	// Limpar dados de teste
	database.DB.Unscoped().Delete(&response)
}

func TestUserLogin(t *testing.T) {
	// Setup
	helpers.SetupTestDB()
	defer helpers.CleanupTestDB()

	// Criar usuário de teste
	user := helpers.CreateTestUser()

	// Criar payload de login
	loginPayload := models.LoginRequest{
		Email:    user.Email,
		Password: "Test123!",
	}

	payload, _ := json.Marshal(loginPayload)

	// Criar request
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	// Criar response recorder
	w := httptest.NewRecorder()

	// Executar request
	router := routes.SetupRoutes()
	router.ServeHTTP(w, req)

	// Verificar status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Verificar resposta
	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token")
	assert.NotEmpty(t, response["token"])
}
