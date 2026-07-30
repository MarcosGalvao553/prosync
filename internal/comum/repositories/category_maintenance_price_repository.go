package repositories

import (
	"database/sql"
	"fmt"
	"strconv"

	"prosync/internal/comum/models"
)

// CategoryMaintenancePriceRepository gerencia preços de manutenção por categoria
type CategoryMaintenancePriceRepository struct {
	db *sql.DB
}

// NovoCategoryMaintenancePriceRepository cria novo repositório
func NovoCategoryMaintenancePriceRepository(db *sql.DB) *CategoryMaintenancePriceRepository {
	return &CategoryMaintenancePriceRepository{db: db}
}

// BuscarPorCategoriaID busca o preço de manutenção específico de uma categoria
func (r *CategoryMaintenancePriceRepository) BuscarPorCategoriaID(categoryID int) (*models.CategoryMaintenancePrice, error) {
	query := `
		SELECT id, category_id, maintenance_price, created_at, updated_at
		FROM category_maintenance_price
		WHERE category_id = ?
		LIMIT 1
	`

	var cmp models.CategoryMaintenancePrice
	err := r.db.QueryRow(query, categoryID).Scan(
		&cmp.ID, &cmp.CategoryID, &cmp.MaintenancePrice, &cmp.CreatedAt, &cmp.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar maintenance price da categoria %d: %w", categoryID, err)
	}

	return &cmp, nil
}

// BuscarPrecoDefault busca o preço de manutenção padrão em system_config_param_values
func (r *CategoryMaintenancePriceRepository) BuscarPrecoDefault() (float64, error) {
	query := `
		SELECT value
		FROM system_config_param_values
		WHERE code = 'maintenance_price'
		LIMIT 1
	`

	var valueStr string
	err := r.db.QueryRow(query).Scan(&valueStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return 3.50, nil // fallback se não encontrar no banco
		}
		return 3.50, fmt.Errorf("erro ao buscar maintenance_price padrão: %w", err)
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 3.50, fmt.Errorf("valor inválido para maintenance_price '%s': %w", valueStr, err)
	}

	return value, nil
}

// BuscarPrecoEfetivo retorna o preço de manutenção da categoria, ou o padrão se não houver específico
func (r *CategoryMaintenancePriceRepository) BuscarPrecoEfetivo(categoryID int) (float64, error) {
	if categoryID > 0 {
		cmp, err := r.BuscarPorCategoriaID(categoryID)
		if err != nil {
			return 0, err
		}
		if cmp != nil {
			return cmp.MaintenancePrice, nil
		}
	}

	return r.BuscarPrecoDefault()
}
