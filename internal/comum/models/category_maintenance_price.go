package models

import "time"

// CategoryMaintenancePrice armazena um preço de manutenção específico por categoria
type CategoryMaintenancePrice struct {
	ID               int       `db:"id" json:"id"`
	CategoryID       int       `db:"category_id" json:"category_id"`
	MaintenancePrice float64   `db:"maintenance_price" json:"maintenance_price"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

// TableName retorna o nome da tabela no banco de dados
func (CategoryMaintenancePrice) TableName() string {
	return "category_maintenance_price"
}
