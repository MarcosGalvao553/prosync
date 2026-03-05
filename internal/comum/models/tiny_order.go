package models

import (
	"database/sql"
	"time"
)

// TinyOrder representa um pedido do Tiny
type TinyOrder struct {
	ID              int            `db:"id" json:"id"`
	ShippingOrderID sql.NullInt64  `db:"shipping_order_id" json:"shipping_order_id,omitempty"`
	OrderTinyID     sql.NullString `db:"order_tiny_id" json:"order_tiny_id,omitempty"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at"`
}

// TableName retorna o nome da tabela no banco de dados
func (TinyOrder) TableName() string {
	return "tiny_orders"
}
