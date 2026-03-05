package repositories

import (
	"database/sql"
	"fmt"

	"prosync/internal/comum/models"
)

// ProductUserRepository gerencia operações com product_users
type ProductUserRepository struct {
	db *sql.DB
}

// NovoProductUserRepository cria novo repositório
func NovoProductUserRepository(db *sql.DB) *ProductUserRepository {
	return &ProductUserRepository{db: db}
}

// ListarPorProductID retorna todos os usuários vinculados a um produto
func (r *ProductUserRepository) ListarPorProductID(productID int) ([]models.ProductUser, error) {
	query := `
		SELECT id, user_id, product_id, tiny_product_id, bling_product_id, 
		       price, created_at, updated_at
		FROM product_users 
		WHERE product_id = ?
		AND user_id in(739,34)
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar product_users: %w", err)
	}
	defer rows.Close()

	var productUsers []models.ProductUser
	for rows.Next() {
		var pu models.ProductUser
		err := rows.Scan(
			&pu.ID, &pu.UserID, &pu.ProductID, &pu.TinyProductID,
			&pu.BlingProductID, &pu.Price, &pu.CreatedAt, &pu.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao fazer scan de product_user: %w", err)
		}
		productUsers = append(productUsers, pu)
	}

	return productUsers, nil
}

// AtualizarBlingProductID atualiza o bling_product_id de um registro
func (r *ProductUserRepository) AtualizarBlingProductID(id int, blingProductID string) error {
	query := `
		UPDATE product_users 
		SET bling_product_id = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.Exec(query, blingProductID, id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar bling_product_id: %w", err)
	}

	return nil
}

// BuscarPorUserIDEProductID busca um registro específico de vínculo
func (r *ProductUserRepository) BuscarPorUserIDEProductID(userID, productID int) (*models.ProductUser, error) {
	query := `
		SELECT id, user_id, product_id, tiny_product_id, bling_product_id, 
		       price, created_at, updated_at
		FROM product_users 
		WHERE user_id = ? AND product_id = ?
		LIMIT 1
	`

	var pu models.ProductUser
	err := r.db.QueryRow(query, userID, productID).Scan(
		&pu.ID, &pu.UserID, &pu.ProductID, &pu.TinyProductID,
		&pu.BlingProductID, &pu.Price, &pu.CreatedAt, &pu.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar product_user: %w", err)
	}

	return &pu, nil
}
