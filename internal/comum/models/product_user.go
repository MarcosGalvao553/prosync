package models

import (
	"database/sql"
)

// ProductUser representa a relação entre produto e usuário no banco de dados
type ProductUser struct {
	ID             int             `db:"id" json:"id"`
	UserID         int             `db:"user_id" json:"user_id"`
	ProductID      int             `db:"product_id" json:"product_id"`
	TinyProductID  sql.NullString  `db:"tiny_product_id" json:"tiny_product_id,omitempty"`
	BlingProductID sql.NullString  `db:"bling_product_id" json:"bling_product_id,omitempty"`
	Price          sql.NullFloat64 `db:"price" json:"price,omitempty"`
	CreatedAt      sql.NullTime    `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt      sql.NullTime    `db:"updated_at" json:"updated_at,omitempty"`
}
