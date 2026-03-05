package models

import (
	"database/sql"
	"time"
)

// Product representa um produto no sistema
type Product struct {
	ID            int             `db:"id" json:"id"`
	Name          string          `db:"name" json:"name"`
	Description   sql.NullString  `db:"description" json:"description,omitempty"`
	Price         sql.NullFloat64 `db:"price" json:"price,omitempty"`
	CostPrice     float64         `db:"cost_price" json:"cost_price"`
	CategoryID    int             `db:"category_id" json:"category_id"`
	IsEnabled     bool            `db:"isEnabled" json:"isEnabled"`
	IsPreSale     bool            `db:"isPreSale" json:"isPreSale"`
	SaleCount     sql.NullInt64   `db:"sale_count" json:"sale_count,omitempty"`
	ReviewCount   sql.NullInt64   `db:"review_count" json:"review_count,omitempty"`
	SKU           sql.NullString  `db:"sku" json:"sku,omitempty"`
	Stock         sql.NullInt64   `db:"stock" json:"stock,omitempty"`
	Observation   sql.NullString  `db:"observation" json:"observation,omitempty"`
	Weight        sql.NullFloat64 `db:"weight" json:"weight,omitempty"`
	Height        sql.NullFloat64 `db:"height" json:"height,omitempty"`
	Width         sql.NullFloat64 `db:"width" json:"width,omitempty"`
	Length        sql.NullFloat64 `db:"length" json:"length,omitempty"`
	ProductTiny   sql.NullString  `db:"product_tiny" json:"product_tiny,omitempty"`
	NCM           sql.NullString  `db:"ncm" json:"ncm,omitempty"`
	EAN           sql.NullString  `db:"ean" json:"ean,omitempty"`
	Marca         sql.NullString  `db:"marca" json:"marca,omitempty"`
	CEST          sql.NullString  `db:"cest" json:"cest,omitempty"`
	StopStock     sql.NullInt64   `db:"stop_stock" json:"stop_stock,omitempty"`
	PromotionID   sql.NullInt64   `db:"promotion_id" json:"promotion_id,omitempty"`
	OriginalPrice sql.NullFloat64 `db:"original_price" json:"original_price,omitempty"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at" json:"updated_at"`
}

// TableName retorna o nome da tabela no banco de dados
func (Product) TableName() string {
	return "products"
}
