package main

import (
	"fmt"
	"github.com/jeffemart/Gotham/settings"
)

func main() {
	// Carregar configurações
	config := settings.LoadSettings()

	// Exemplo de uso
	fmt.Println("Conectando ao banco de dados:")
	fmt.Printf("Driver: %s\nHost: %s\nPort: %s\nUser: %s\n",
		config.Database.Driver,
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
	)

	// A partir daqui, você pode usar essas configurações para conectar ao banco, Redis, etc.
}
