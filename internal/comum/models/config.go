package models

import "database/sql"

// Config representa uma configuração geral no banco de dados
type Config struct {
	ID          int            `db:"id" json:"id"`
	Code        string         `db:"code" json:"code"`
	Description sql.NullString `db:"description" json:"description,omitempty"`
	Value       sql.NullString `db:"value" json:"value,omitempty"`
}
