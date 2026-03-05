package models

import (
	"database/sql"
	"time"
)

// SystemConfigParam representa um parâmetro de configuração do sistema
type SystemConfigParam struct {
	ID             int            `db:"id" json:"id"`
	Name           sql.NullString `db:"name" json:"name,omitempty"`
	Description    sql.NullString `db:"description" json:"description,omitempty"`
	ShowToUser     sql.NullBool   `db:"show_to_user" json:"show_to_user,omitempty"`
	Code           sql.NullString `db:"code" json:"code,omitempty"`
	SystemConfigID sql.NullInt64  `db:"system_config_id" json:"system_config_id,omitempty"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
}

// TableName retorna o nome da tabela no banco de dados
func (SystemConfigParam) TableName() string {
	return "system_config_params"
}
