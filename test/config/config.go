package config

import (
	"os"
)

func init() {
	// Configurações para ambiente de teste
	os.Setenv("APP_ENV", "test")
	os.Setenv("APP_KEY", "test_key_for_jwt_signing")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "gotham_test")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
}
