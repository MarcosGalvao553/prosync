package models

import (
	"database/sql"
	"time"
)

// SystemConfigParamValue representa um valor de parâmetro de configuração
type SystemConfigParamValue struct {
	ID                  int            `db:"id" json:"id"`
	Name                sql.NullString `db:"name" json:"name,omitempty"`
	Code                sql.NullString `db:"code" json:"code,omitempty"`
	Value               sql.NullString `db:"value" json:"value,omitempty"`
	UserID              sql.NullInt64  `db:"user_id" json:"user_id,omitempty"`
	SystemConfigParamID sql.NullInt64  `db:"system_config_param_id" json:"system_config_param_id,omitempty"`
	IsConfig            sql.NullBool   `db:"is_config" json:"is_config,omitempty"`
	CreatedAt           time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at" json:"updated_at"`
}

// TableName retorna o nome da tabela no banco de dados
func (SystemConfigParamValue) TableName() string {
	return "system_config_param_values"
}
