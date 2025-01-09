package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jeffemart/Gotham/models"
)

// Define a chave secreta (pode ser mais segura, por exemplo, em variáveis de ambiente)
var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Claims personalizados para incluir a role do usuário
type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

// Função para gerar um novo token JWT
func GenerateToken(user models.User) (string, error) {
	claims := Claims{
		Email: user.Email,
		Role:  user.Role,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "GothamApp",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Função para fazer o parse e validação do token JWT
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	// Parse do token com as claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar o método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token inválido: método de assinatura inválido")
		}
		// Retornar a chave secreta
		return secretKey, nil
	})

	// Verificar se houve erro ao parsear
	if err != nil {
		return nil, nil, err
	}

	// Extrair as claims do token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return token, claims, nil
	} else {
		return nil, nil, fmt.Errorf("token inválido")
	}
}

// Função para validar o token
func ValidateToken(tokenString string) (*Claims, error) {
	// Faz o parse do token
	_, claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("token inválido: %v", err)
	}

	// Verifica se o token expirou
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token expirado")
	}

	// Retorna as claims do token
	return claims, nil
}

// Função para verificar se o token expirou
func TokenExpired(token *jwt.Token) bool {
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return true
	}
	return claims.ExpiresAt < time.Now().Unix()
}
