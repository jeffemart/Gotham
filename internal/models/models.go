package models

import "time"

// User representa um usuário no banco de dados
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:255;not null;unique"`
	Password  string    `gorm:"size:255;not null"`
	RoleID    uint      `gorm:"not null"`
	Role      Role      `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Role struct {
	ID          uint         `gorm:"primaryKey"`
	Name        string       `gorm:"size:255;not null;unique"`
	Permissions []Permission `gorm:"many2many:role_permissions"`
	Capabilities []string    `gorm:"type:text[]"`
}

type Permission struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255;not null;unique"`
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// PaginatedResponse representa uma resposta paginada
type PaginatedResponse struct {
	Status      int         `json:"status"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data"`
	TotalCount  int         `json:"total_count"`
	TotalPages  int         `json:"total_pages"`
	CurrentPage int         `json:"current_page"`
	Limit       int         `json:"limit"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Adicionar constantes para capacidades
const (
	CapabilityCreateUser    = "create:user"
	CapabilityReadUser      = "read:user"
	CapabilityUpdateUser    = "update:user"
	CapabilityDeleteUser    = "delete:user"
	CapabilityManageRoles   = "manage:roles"
	CapabilityViewTasks     = "view:tasks"
	CapabilityManageTasks   = "manage:tasks"
)
