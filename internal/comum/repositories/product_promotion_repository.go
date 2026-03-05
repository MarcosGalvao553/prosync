package repositories

import (
	"database/sql"
)

type ProductPromotionRepository struct {
	db *sql.DB
}

func NovoProductPromotionRepository(db *sql.DB) *ProductPromotionRepository {
	return &ProductPromotionRepository{db: db}
}

// VerificarPromocaoAtiva verifica se o produto tem promoção ativa
func (r *ProductPromotionRepository) VerificarPromocaoAtiva(productID int) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM product_promotions 
		WHERE product_id = ? 
		AND active = 1
		AND deprecated = 0
		AND start_date <= NOW()
		AND end_date >= NOW()
	`

	var count int
	err := r.db.QueryRow(query, productID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
