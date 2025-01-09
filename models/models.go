package models

import "time"

// User representa um usuário no banco de dados
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:255;not null;unique"`
	Password  string    `gorm:"size:255;not null"`
	Role      string    `gorm:"size:20;not null;check:role IN ('cliente', 'agente', 'admin')"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}