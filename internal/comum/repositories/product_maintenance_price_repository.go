package repositories

import (
	"database/sql"
	"fmt"

	"prosync/internal/comum/models"
)

// ProductMaintenancePriceRepository gerencia preços de manutenção por produto
type ProductMaintenancePriceRepository struct {
	db           *sql.DB
	categoryRepo *CategoryMaintenancePriceRepository
}

// NovoProductMaintenancePriceRepository cria novo repositório
func NovoProductMaintenancePriceRepository(db *sql.DB, categoryRepo *CategoryMaintenancePriceRepository) *ProductMaintenancePriceRepository {
	return &ProductMaintenancePriceRepository{db: db, categoryRepo: categoryRepo}
}

// BuscarPorProdutoID busca o preço de manutenção específico de um produto
func (r *ProductMaintenancePriceRepository) BuscarPorProdutoID(productID int) (*models.ProductMaintenancePrice, error) {
	query := `
		SELECT id, product_id, maintenance_price, created_at, updated_at
		FROM product_maintenance_price
		WHERE product_id = ?
		LIMIT 1
	`

	var pmp models.ProductMaintenancePrice
	err := r.db.QueryRow(query, productID).Scan(
		&pmp.ID, &pmp.ProductID, &pmp.MaintenancePrice, &pmp.CreatedAt, &pmp.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar maintenance price do produto %d: %w", productID, err)
	}

	return &pmp, nil
}

// BuscarPrecoEfetivo retorna o preço de manutenção seguindo a ordem de prioridade:
// 1. Produto específico
// 2. Categoria do produto
// 3. Configuração global
func (r *ProductMaintenancePriceRepository) BuscarPrecoEfetivo(productID, categoryID int) (float64, error) {
	if productID > 0 {
		pmp, err := r.BuscarPorProdutoID(productID)
		if err != nil {
			return 0, err
		}
		if pmp != nil {
			return pmp.MaintenancePrice, nil
		}
	}

	return r.categoryRepo.BuscarPrecoEfetivo(categoryID)
}
