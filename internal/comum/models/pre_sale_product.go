package models

import (
	"database/sql"
	"time"
)

// PreSaleProduct representa um produto em pré-venda
type PreSaleProduct struct {
	ID        int          `db:"id" json:"id"`
	ProductID int          `db:"product_id" json:"product_id"`
	EndDate   sql.NullTime `db:"end_date" json:"end_date,omitempty"`
	Active    bool         `db:"active" json:"active"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
}

// TableName retorna o nome da tabela no banco de dados
func (PreSaleProduct) TableName() string {
	return "pre_sale_products"
}
