package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jeffemart/Gotham/internal/database"
	"github.com/jeffemart/Gotham/internal/models"
	"github.com/jeffemart/Gotham/internal/utils"
)

// AuthMiddleware verifica se o token JWT é válido
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar token JWT no cabeçalho Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token não fornecido", http.StatusUnauthorized)
			return
		}

		// Extrair o token do cabeçalho
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse e validação do token
		token, claims, err := utils.ParseToken(tokenString)
		if err != nil || token == nil || !token.Valid {
			http.Error(w, "Token inválido ou expirado", http.StatusUnauthorized)
			return
		}

		// Adicionar o Claims ao contexto para que as próximas funções possam acessar
		ctx := context.WithValue(r.Context(), utils.RoleKey, claims)
		r = r.WithContext(ctx)

		// Chamar o próximo handler na cadeia
		next.ServeHTTP(w, r)
	})
}

// RoleMiddleware verifica se o usuário tem uma das roles necessárias para acessar a rota
func RoleMiddleware(roles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obter claims do contexto (setado pelo AuthMiddleware)
			claims, ok := r.Context().Value(utils.RoleKey).(*utils.Claims)
			if !ok {
				http.Error(w, "Erro ao obter informações do usuário", http.StatusInternalServerError)
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
			default:
				// Para outros métodos, concede permissão
				permissionGranted = true
			}

			// Se a permissão não for concedida, retorna erro
			if !permissionGranted {
				http.Error(w, "Acesso negado: você não tem permissão para acessar essa rota", http.StatusForbidden)
				return
			}

			// Chamar o próximo handler na cadeia
			next.ServeHTTP(w, r)
		})
	}
}

// CapabilityMiddleware verifica se o usuário tem as capacidades necessárias
func CapabilityMiddleware(requiredCapabilities ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obter claims do contexto
			claims, ok := r.Context().Value(utils.RoleKey).(*utils.Claims)
			if !ok {
				http.Error(w, "Erro ao obter informações do usuário", http.StatusInternalServerError)
				return
			}

			// Recuperar a role do banco de dados
			var role models.Role
			if err := database.DB.First(&role, claims.RoleID).Error; err != nil {
				http.Error(w, "Role não encontrada", http.StatusForbidden)
				return
			}

			// Verificar se o usuário tem todas as capacidades necessárias
			hasAllCapabilities := true
			for _, requiredCap := range requiredCapabilities {
				hasCapability := false
				for _, userCap := range role.Capabilities {
					if userCap == requiredCap || userCap == "*" { // "*" representa acesso total
						hasCapability = true
						break
					}
				}
				if !hasCapability {
					hasAllCapabilities = false
					break
				}
			}

			if !hasAllCapabilities {
				http.Error(w, "Acesso negado: capacidades insuficientes", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
