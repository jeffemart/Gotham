package unit

import (
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jeffemart/Gotham/internal/models"
	"github.com/jeffemart/Gotham/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestTokenGeneration(t *testing.T) {
	// Criar usuário de teste
	user := models.User{
		ID:    1,
		Email: "test@example.com",
		Role: models.Role{
			ID:   1,
			Name: "admin",
			Permissions: []models.Permission{
				{ID: 1, Name: "read"},
				{ID: 2, Name: "write"},
			},
		},
	}

	// Gerar token
	token, err := utils.GenerateTokenWithPermissions(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validar token
	claims, err := utils.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, uint(1), claims.RoleID)
	assert.Contains(t, claims.Permissions, "read")
	assert.Contains(t, claims.Permissions, "write")
}

func TestTokenExpiration(t *testing.T) {
	// Criar token expirado
	claims := &utils.Claims{
		Email:  "test@example.com",
		RoleID: 1,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(os.Getenv("APP_KEY")))

	// Validar token expirado
	_, err := utils.ValidateToken(tokenString)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token expirado")
}
