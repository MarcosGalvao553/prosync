package models

import "database/sql"

// ProductImage representa uma imagem de produto
type ProductImage struct {
	ID            int            `db:"id" json:"id"`
	ImageType     int            `db:"image_type" json:"image_type"`
	ImageSrc      string         `db:"image_src" json:"image_src"`
	ProductID     int            `db:"product_id" json:"product_id"`
	ImageSrcSmall sql.NullString `db:"Image_src_small" json:"Image_src_small,omitempty"`
}

// TableName retorna o nome da tabela no banco de dados
func (ProductImage) TableName() string {
	return "product_images"
}
