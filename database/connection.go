package database

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jeffemart/Gotham/settings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Definindo o contexto global para ser usado nas operações Redis
var Ctx = context.Background()

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

// Connect estabelece a conexão com o banco de dados e o Redis
func Connect() {
	// Carregar as configurações das variáveis de ambiente
	config := settings.LoadSettings()

	// Conexão com PostgreSQL
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.Host,
		config.Database.Port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	log.Println("Conexão com o banco de dados PostgreSQL estabelecida.")

	// Conexão com Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password, // Senha (se configurada)
		DB:       config.Redis.DB,       // Banco (default: 0)
	})

	// Testar conexão com Redis
	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Erro ao conectar ao Redis: %v", err)
	}
	log.Println("Conexão com o Redis estabelecida.")
}


// package database

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/go-redis/redis/v8"
// 	"github.com/jeffemart/Gotham/settings"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// var (
// 	DB          *gorm.DB
// 	RedisClient *redis.Client
// )

// // Connect estabelece a conexão com o banco de dados e o Redis
// func Connect() {
// 	// Carregar as configurações das variáveis de ambiente
// 	config := settings.LoadSettings()

// 	// Conexão com PostgreSQL
// 	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
// 		config.Database.User,
// 		config.Database.Password,
// 		config.Database.Name,
// 		config.Database.Host,
// 		config.Database.Port,
// 	)

// 	var err error
// 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
// 	}
// 	log.Println("Conexão com o banco de dados PostgreSQL estabelecida.")

// 	// Conexão com Redis
// 	RedisClient = redis.NewClient(&redis.Options{
// 		Addr:     fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
// 		Password: config.Redis.Password, // Senha (se configurada)
// 		DB:       config.Redis.DB,       // Banco (default: 0)
// 	})

// 	// Testar conexão com Redis
// 	_, err = RedisClient.Ping(context.Background()).Result()
// 	if err != nil {
// 		log.Fatalf("Erro ao conectar ao Redis: %v", err)
// 	}
// 	log.Println("Conexão com o Redis estabelecida.")
// }
