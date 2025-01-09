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
	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnv("DB_PORT", "5432")
	config.Database.User = getEnv("DB_USER", "myuser")
	config.Database.Password = getEnv("DB_PASSWORD", "mypassword")
	config.Database.Name = getEnv("DB_NAME", "gotham_db")

	// Configurações do Redis
	config.Redis.Host = getEnv("REDIS_HOST", "localhost")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
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
