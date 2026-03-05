package models

import (
	"database/sql"
	"time"
)

// ProductPromotion representa uma promoção de produto
type ProductPromotion struct {
	ID                   int            `db:"id" json:"id"`
	ProductID            int            `db:"product_id" json:"product_id"`
	PromotionalPrice     float64        `db:"promotional_price" json:"promotional_price"`
	ProductOriginalPrice float64        `db:"product_original_price" json:"product_original_price"`
	Active               bool           `db:"active" json:"active"`
	UserID               sql.NullInt64  `db:"user_id" json:"user_id,omitempty"`
	StartDate            time.Time      `db:"start_date" json:"start_date"`
	EndDate              time.Time      `db:"end_date" json:"end_date"`
	Description          sql.NullString `db:"description" json:"description,omitempty"`
	Deprecated           bool           `db:"deprecated" json:"deprecated"`
	CreatedAt            time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time      `db:"updated_at" json:"updated_at"`
}

// TableName retorna o nome da tabela no banco de dados
func (ProductPromotion) TableName() string {
	return "product_promotions"
}
