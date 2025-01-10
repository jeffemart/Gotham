package settings

import (
	"fmt"
	"log"
	"os"
)

// Config estrutura que contém as configurações do projeto
type Config struct {
	Database struct {
		Driver   string // Driver do banco de dados (e.g., postgres, mysql, mongo)
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
		DB       int
	}
}

// LoadSettings carrega as configurações das variáveis de ambiente
func LoadSettings() *Config {
	config := &Config{}

	// Configurações do banco de dados
	config.Database.Driver = getEnv("DB_DRIVER", "postgres")
	config.Database.Host = getEnv("DB_HOST", "172.28.0.3") // IP do Postgres
	config.Database.Port = getEnv("POSTGRES_PORT", "5432")
	config.Database.User = getEnv("POSTGRES_USER", "user")
	config.Database.Password = getEnv("POSTGRES_PASSWORD", "password")
	config.Database.Name = getEnv("POSTGRES_DB", "gotham_db")

	// Configurações do Redis
	config.Redis.Host = getEnv("REDIS_HOST", "172.28.0.4")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "password")
	config.Redis.DB = getEnvAsInt("REDIS_DB", 0)

	return config
}

// Função auxiliar para obter variáveis de ambiente
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Função auxiliar para obter variáveis de ambiente como inteiro
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	var value int
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		log.Printf("Erro ao converter %s para inteiro: %v. Usando valor padrão: %d", key, err, defaultValue)
		return defaultValue
	}
	return value
}
