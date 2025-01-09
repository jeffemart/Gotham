package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jeffemart/Gotham/utils"
)

// Defina um tipo customizado para o tipo string
type RoleKey string

// RoleMiddleware verifica se o usuário tem uma das roles permitidas
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

			// Verificar se a role do usuário está na lista de roles permitidas
			roleAllowed := false
			for _, role := range roles {
				if role == claims.Role {
					roleAllowed = true
					break
				}
			}

			// Se a role não for permitida, retorna erro
			if !roleAllowed {
				http.Error(w, "Acesso negado: você não tem permissão para acessar essa rota", http.StatusForbidden)
				return
			}

			// Adicionar o usuário ao contexto para que as próximas funções possam acessar
			ctx := context.WithValue(r.Context(), RoleKey("user"), claims)
			r = r.WithContext(ctx)

			// Chamar o próximo handler na cadeia
			next.ServeHTTP(w, r)
		})
	}
}
