package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/jeffemart/Gotham/models"
)

// Contexto para uso com Redis
var ctx = context.Background()

// Inicializar cliente Redis
var RedisClient = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	Password: os.Getenv("REDIS_PASSWORD"), // Deixe vazio se não houver senha
	DB:       0,                          // Banco padrão
})

// Define a chave secreta (deve ser configurada como variável de ambiente)
var secretKey = []byte(os.Getenv("APP_KEY"))

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
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	// Armazena o token no Redis com o mesmo tempo de expiração
	err = RedisClient.Set(ctx, tokenString, "valid", 24*time.Hour).Err()
	if err != nil {
		return "", fmt.Errorf("erro ao salvar token no Redis: %v", err)
	}

	return tokenString, nil
}

// Função para fazer o parse e validação do token JWT
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token inválido: método de assinatura inválido")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		return token, claims, nil
	}

	return nil, nil, fmt.Errorf("token inválido")
}

// Função para validar o token, incluindo verificação no Redis
func ValidateToken(tokenString string) (*Claims, error) {
	_, claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("token inválido: %v", err)
	}

	// Verifica se o token expirou
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token expirado")
	}

	// Verifica no Redis se o token ainda é válido
	redisValue, err := RedisClient.Get(ctx, tokenString).Result()
	if err == redis.Nil || redisValue != "valid" {
		return nil, fmt.Errorf("token revogado ou inválido")
	}

	return claims, nil
}

// Função para revogar o token, removendo-o do Redis
func RevokeToken(tokenString string) error {
	err := RedisClient.Del(ctx, tokenString).Err()
	if err != nil {
		return fmt.Errorf("erro ao revogar token: %v", err)
	}
	return nil
}

// Função para verificar se o token expirou (para uso interno)
func TokenExpired(token *jwt.Token) bool {
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return true
	}
	return claims.ExpiresAt < time.Now().Unix()
}
