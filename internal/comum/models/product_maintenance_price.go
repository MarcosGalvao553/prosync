package models

import "time"

// ProductMaintenancePrice armazena um preço de manutenção específico por produto
type ProductMaintenancePrice struct {
	ID               int       `db:"id" json:"id"`
	ProductID        int       `db:"product_id" json:"product_id"`
	MaintenancePrice float64   `db:"maintenance_price" json:"maintenance_price"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

// TableName retorna o nome da tabela no banco de dados
func (ProductMaintenancePrice) TableName() string {
	return "product_maintenance_price"
}
