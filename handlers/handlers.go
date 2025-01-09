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
func Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	var user models.User
	result := database.DB.Where("email = ?", loginRequest.Email).First(&user)
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

	// Gera o token JWT
	token, err := utils.GenerateToken(user)
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

	token, err := utils.GenerateToken(user)
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

// GetUsers retorna todos os usuários cadastrados
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	result := database.DB.Find(&users)
	if result.Error != nil {
		http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// UpdateUser atualiza as informações de um usuário
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

	user.ID = uint(id)
	user.UpdatedAt = time.Now()
	result := database.DB.Save(&user)
	if result.Error != nil {
		http.Error(w, "Erro ao atualizar usuário", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// DeleteUser remove um usuário pelo ID
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
