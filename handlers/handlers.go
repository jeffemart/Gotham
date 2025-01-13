package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jeffemart/Gotham/database"
	"github.com/jeffemart/Gotham/models"
	"github.com/jeffemart/Gotham/utils"
	"golang.org/x/crypto/bcrypt"
)

// Login autentica o usuário e gera um token JWT
// @Summary Login do usuário
// @Description Autentica o usuário e gera um token JWT
// @Accept  json
// @Produce  json
// @Param loginRequest body models.LoginRequest true "Credenciais do usuário"
// @Success 200 {object} map[string]string "Token gerado"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 401 {string} string "Usuário não encontrado ou senha incorreta"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	var user models.User
	result := database.DB.Where("email = ?", loginRequest.Email).Preload("Role.Permissions").First(&user)
	if result.Error != nil {
		http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	}

	// Verifica a senha
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		http.Error(w, "Senha incorreta", http.StatusUnauthorized)
		return
	}

	// Gera o token JWT com as permissões do papel do usuário
	token, err := utils.GenerateTokenWithPermissions(user)
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}

	// Retorna o token
	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RefreshToken gera um novo token JWT usando o refresh token
// @Summary Gera um novo token JWT usando o refresh token
// @Description Gera um novo token JWT quando fornecido um refresh token válido
// @Accept  json
// @Produce  json
// @Param refreshRequest body models.RefreshRequest true "Refresh Token"
// @Success 200 {object} map[string]string "Novo token gerado"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 401 {string} string "Refresh token inválido"
// @Failure 404 {string} string "Usuário não encontrado"
// @Router /refresh_token [post]
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshRequest models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshRequest); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Valida o refresh token
	claims, err := utils.ValidateToken(refreshRequest.RefreshToken)
	if err != nil {
		http.Error(w, "Refresh token inválido", http.StatusUnauthorized)
		return
	}

	// Gera um novo token usando o email da claim
	var user models.User
	result := database.DB.Where("email = ?", claims.Email).First(&user)
	if result.Error != nil {
		http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateTokenWithPermissions(user)
	if err != nil {
		http.Error(w, "Erro ao gerar novo token", http.StatusInternalServerError)
		return
	}

	// Retorna o novo token
	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateUser cria um novo usuário
// @Summary Cria um novo usuário
// @Description Cria um novo usuário no sistema com os dados fornecidos
// @Accept  json
// @Produce  json
// @Param user body models.User true "Dados do novo usuário"
// @Success 201 {object} models.User "Usuário criado com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro ao criar usuário"
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// Decodifica o corpo da requisição para o objeto user
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Criptografa a senha do usuário com bcrypt antes de salvar no banco
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erro ao criptografar senha", http.StatusInternalServerError)
		return
	}

	// Substitui a senha do usuário pela versão criptografada
	user.Password = string(hashedPassword)

	// Salva o usuário no banco de dados
	result := database.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Erro ao criar usuário", http.StatusInternalServerError)
		return
	}

	// Retorna o usuário criado com status 201 (Created)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser retorna um usuário pelo ID
// @Summary Retorna um usuário pelo ID
// @Description Obtém um usuário específico com base no ID fornecido
// @Produce  json
// @Param id path int true "ID do usuário"
// @Success 200 {object} models.User "Usuário encontrado"
// @Failure 400 {string} string "ID inválido"
// @Failure 404 {string} string "Usuário não encontrado"
// @Router /users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var user models.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// GetUsers retorna todos os usuários cadastrados com paginação
// @Summary Retorna todos os usuários com paginação
// @Description Obtém uma lista de todos os usuários cadastrados no sistema com base na paginação
// @Produce  json
// @Param page query int false "Número da página" default(1)
// @Param limit query int false "Limite de usuários por página" default(10)
// @Success 200 {object} utils.PaginatedResponse "Lista de usuários com paginação"
// @Failure 500 {string} string "Erro ao buscar usuários"
// @Router /users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Pega os parâmetros de query (página e limite)
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Definir página e limite com valores padrão
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	limit := 10
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
		if limit < 1 {
			limit = 10
		}
	}

	// Calcula o deslocamento (offset) baseado na página e no limite
	offset := (page - 1) * limit

	// Contagem total de usuários
	var totalCount int64
	if err := database.DB.Model(&models.User{}).Count(&totalCount).Error; err != nil {
		http.Error(w, "Erro ao contar usuários", http.StatusInternalServerError)
		return
	}

	// Consultar usuários com limite e offset para paginação
	var users []models.User
	if err := database.DB.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
		return
	}

	// Calcula o número total de páginas
	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	// Monta a resposta paginada
	response := utils.PaginatedResponse{
		Status:      http.StatusOK,
		Message:     "Usuários encontrados",
		Data:        users,
		TotalCount:  int(totalCount),
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
	}

	// Retorna a resposta como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateUser atualiza as informações de um usuário
// @Summary Atualiza as informações de um usuário
// @Description Atualiza os dados de um usuário com base no ID fornecido
// @Accept  json
// @Security BearerAuth
// @Produce  json
// @Param id path int true "ID do usuário"
// @Param user body models.User true "Dados do usuário a serem atualizados"
// @Success 200 {object} models.User "Usuário atualizado com sucesso"
// @Failure 400 {string} string "ID ou dados inválidos"
// @Failure 404 {string} string "Usuário não encontrado"
// @Failure 500 {string} string "Erro ao atualizar usuário"
// @Router /admin/users/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Carrega o usuário do banco de dados
	var existingUser models.User
	result := database.DB.First(&existingUser, id)
	if result.Error != nil {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	// Atualiza os campos fornecidos na requisição, se não forem os valores padrão (zero)
	if user.Email != "" && user.Email != existingUser.Email {
		existingUser.Email = user.Email
	}
	if user.Name != "" && user.Name != existingUser.Name {
		existingUser.Name = user.Name
	}
	if user.Password != "" { // Criptografa a senha se fornecida
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Erro ao criptografar senha", http.StatusInternalServerError)
			return
		}
		existingUser.Password = string(hashedPassword)
	}
	if user.RoleID != 0 && user.RoleID != existingUser.RoleID {
		existingUser.RoleID = user.RoleID
	}

	// Atualiza os campos no banco de dados
	existingUser.UpdatedAt = time.Now()
	result = database.DB.Save(&existingUser)
	if result.Error != nil {
		http.Error(w, "Erro ao atualizar usuário", http.StatusInternalServerError)
		return
	}

	// Retorna o usuário atualizado
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingUser)
}

// DeleteUser remove um usuário pelo ID
// @Summary Remove um usuário pelo ID
// @Description Exclui um usuário com base no ID fornecido
// @Security BearerAuth
// @Param id path int true "ID do usuário"
// @Success 200 {object} map[string]string "Usuário excluído com sucesso"
// @Failure 400 {string} string "ID inválido"
// @Failure 500 {string} string "Erro ao excluir usuário"
// @Router /admin/users/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result := database.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		http.Error(w, "Erro ao excluir usuário", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuário excluído com sucesso"})
}

// GetTasks retorna uma lista de tarefas fictícias
// @Summary Retorna uma lista de tarefas
// @Description Obtém uma lista de tarefas fictícias para a demonstração
// @Security BearerAuth
// @Produce  json
// @Success 200 {array} map[string]string "Lista de tarefas"
// @Router /protected/tasks [get]
func GetTasks(w http.ResponseWriter, r *http.Request) {
	// Lista fictícia de tarefas
	tasks := []map[string]string{
		{"id": "1", "task": "Finalizar relatório", "status": "Pendente"},
		{"id": "2", "task": "Enviar e-mail para cliente", "status": "Concluído"},
		{"id": "3", "task": "Atualizar sistema", "status": "Em progresso"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
