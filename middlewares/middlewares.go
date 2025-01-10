package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jeffemart/Gotham/database"
	"github.com/jeffemart/Gotham/models"
	"github.com/jeffemart/Gotham/utils"
)

// RoleMiddleware verifica se o usuário tem uma das permissões ou roles necessárias para acessar a rota
func RoleMiddleware(roles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verificar token JWT no cabeçalho Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Token não fornecido", http.StatusUnauthorized)
				return
			}

			// Extrair o token do cabeçalho
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse do token
			token, claims, err := utils.ParseToken(tokenString)
			if err != nil || token == nil || !token.Valid {
				http.Error(w, "Token inválido ou expirado", http.StatusUnauthorized)
				return
			}

			// Recuperar a role do banco de dados com as permissões associadas
			var role models.Role
			if err := database.DB.Preload("Permissions").First(&role, claims.RoleID).Error; err != nil {
				http.Error(w, "Role não encontrada", http.StatusForbidden)
				return
			}

			// Verificar se a role do usuário corresponde a algum dos papéis fornecidos
			roleFound := false
			for _, allowedRole := range roles {
				if role.Name == allowedRole {
					roleFound = true
					break
				}
			}

			// Se a role não for encontrada, retorna erro de acesso negado
			if !roleFound {
				http.Error(w, "Acesso negado: você não tem a role necessária para acessar essa rota", http.StatusForbidden)
				return
			}

			// Verificar permissões dinâmicas
			permissionGranted := false
			// Para cada permissão da role, verificamos se ela está na lista de permissões associadas
			for _, permission := range role.Permissions {
				// Exemplo de verificação de permissão dinâmica
				// Aqui você deve implementar a lógica necessária para verificar a permissão da rota atual
				if permission.Name == "required_permission_name" { // Substitua pelo nome da permissão exigida
					permissionGranted = true
					break
				}
			}

			// Se a permissão não for concedida, retorna erro
			if !permissionGranted {
				http.Error(w, "Acesso negado: você não tem permissão para acessar essa rota", http.StatusForbidden)
				return
			}

			// Adicionar o Claims ao contexto para que as próximas funções possam acessar
			ctx := context.WithValue(r.Context(), utils.RoleKey, claims)
			r = r.WithContext(ctx)

			// Chamar o próximo handler na cadeia
			next.ServeHTTP(w, r)
		})
	}
}



// package middlewares

// import (
// 	"context"
// 	"net/http"
// 	"strings"

// 	"github.com/gorilla/mux"
// 	"github.com/jeffemart/Gotham/database"
// 	"github.com/jeffemart/Gotham/models"
// 	"github.com/jeffemart/Gotham/utils"
// )

// // RoleMiddleware verifica se o usuário tem uma das permissões ou roles necessárias para acessar a rota
// func RoleMiddleware(roles ...string) mux.MiddlewareFunc {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			// Verificar token JWT no cabeçalho Authorization
// 			authHeader := r.Header.Get("Authorization")
// 			if authHeader == "" {
// 				http.Error(w, "Token não fornecido", http.StatusUnauthorized)
// 				return
// 			}

// 			// Extrair o token do cabeçalho
// 			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

// 			// Parse do token
// 			token, claims, err := utils.ParseToken(tokenString)
// 			if err != nil || token == nil || !token.Valid {
// 				http.Error(w, "Token inválido ou expirado", http.StatusUnauthorized)
// 				return
// 			}

// 			// Recuperar a role do banco de dados com as permissões associadas
// 			var role models.Role
// 			if err := database.DB.Preload("Permissions").First(&role, claims.RoleID).Error; err != nil {
// 				http.Error(w, "Role não encontrada", http.StatusForbidden)
// 				return
// 			}

// 			// Verificar se a role do usuário corresponde a algum dos papéis fornecidos
// 			roleFound := false
// 			for _, allowedRole := range roles {
// 				if role.Name == allowedRole {
// 					roleFound = true
// 					break
// 				}
// 			}

// 			if !roleFound {
// 				http.Error(w, "Acesso negado: você não tem a role necessária para acessar essa rota", http.StatusForbidden)
// 				return
// 			}

// 			// Verificar permissões dinâmicas
// 			permissionGranted := false
// 			for _, permission := range role.Permissions {
// 				// Aqui, você pode verificar a permissão necessária para a rota, por exemplo:
// 				if permission.Name == "required_permission_name" { // Substitua pela lógica necessária
// 					permissionGranted = true
// 					break
// 				}
// 			}

// 			// Se a permissão não for concedida, retorna erro
// 			if !permissionGranted {
// 				http.Error(w, "Acesso negado: você não tem permissão para acessar essa rota", http.StatusForbidden)
// 				return
// 			}

// 			// Adicionar o Claims ao contexto para que as próximas funções possam acessar
// 			ctx := context.WithValue(r.Context(), utils.RoleKey, claims)
// 			r = r.WithContext(ctx)

// 			// Chamar o próximo handler na cadeia
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }
