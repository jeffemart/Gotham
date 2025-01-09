package models

// User representa a tabela users no banco de dados
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(255);not null"`
	Email     string `gorm:"type:varchar(255);not null;unique"`
	Password  string `gorm:"type:varchar(255);not null"`
	Role      string `gorm:"type:varchar(20);not null;check:role in ('cliente', 'agente', 'admin')"`
	CreatedAt string `gorm:"default:current_timestamp"`
	UpdatedAt string `gorm:"default:current_timestamp"`
}
