package models

import (
	"database/sql"
	"time"
)

// SystemConfig representa uma configuração do sistema
type SystemConfig struct {
	ID          int            `db:"id" json:"id"`
	Name        sql.NullString `db:"name" json:"name,omitempty"`
	Description sql.NullString `db:"description" json:"description,omitempty"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}

// TableName retorna o nome da tabela no banco de dados
func (SystemConfig) TableName() string {
	return "system_configs"
}
