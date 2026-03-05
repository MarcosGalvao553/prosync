package models

import "database/sql"

// Category representa uma categoria de produtos
type Category struct {
	ID           int             `db:"id" json:"id"`
	Name         string          `db:"name" json:"name"`
	ImageSrc     sql.NullString  `db:"image_src" json:"image_src,omitempty"`
	Code         sql.NullString  `db:"code" json:"code,omitempty"`
	CategoryID   sql.NullString  `db:"category_id" json:"category_id,omitempty"`
	Range1       sql.NullFloat64 `db:"range_1" json:"range_1,omitempty"`
	Range2       sql.NullFloat64 `db:"range_2" json:"range_2,omitempty"`
	Range3       sql.NullFloat64 `db:"range_3" json:"range_3,omitempty"`
	FreeShipping sql.NullFloat64 `db:"free_shipping" json:"free_shipping,omitempty"`
}

// TableName retorna o nome da tabela no banco de dados
func (Category) TableName() string {
	return "categories"
}
