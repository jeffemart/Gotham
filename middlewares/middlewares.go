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

			// Verificar as permissões dinâmicas baseadas na ação da rota
			permissionGranted := false
			switch r.Method {
			case http.MethodDelete:
				// Para DELETE, precisa de permissão de admin (ID 1)
				if claims.RoleID == 1 {
					permissionGranted = true
				}
			case http.MethodPut, http.MethodPatch:
				// Para PUT e PATCH, precisa de permissão de agent (ID 2) ou admin (ID 1)
				if claims.RoleID == 1 || claims.RoleID == 2 {
					permissionGranted = true
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
