package models

import (
	"database/sql"
	"time"
)

// BlingConfiguration representa a configuração do Bling
type BlingConfiguration struct {
	ID            int            `db:"id" json:"id"`
	ClientID      sql.NullString `db:"client_id" json:"client_id,omitempty"`
	SecretKey     sql.NullString `db:"secret_key" json:"secret_key,omitempty"`
	URLCallback   sql.NullString `db:"url_callback" json:"url_callback,omitempty"`
	Postcode      sql.NullString `db:"postcode" json:"postcode,omitempty"`
	AccessToken   sql.NullString `db:"access_token" json:"access_token,omitempty"`
	RefreshToken  sql.NullString `db:"refresh_token" json:"refresh_token,omitempty"`
	TokenValidate sql.NullTime   `db:"token_validate" json:"token_validate,omitempty"`
	Code          sql.NullString `db:"code" json:"code,omitempty"`
	UserID        sql.NullInt64  `db:"user_id" json:"user_id,omitempty"`
	CreatedAt     time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at" json:"updated_at"`
}

// TableName retorna o nome da tabela no banco de dados
func (BlingConfiguration) TableName() string {
	return "bling_configurations"
}
