package repositories

import (
	"database/sql"
	"fmt"

	"prosync/internal/comum/models"
)

// ProductImageRepository gerencia operações com imagens de produtos
type ProductImageRepository struct {
	db *sql.DB
}

// NovoProductImageRepository cria novo repositório de imagens de produtos
func NovoProductImageRepository(db *sql.DB) *ProductImageRepository {
	return &ProductImageRepository{db: db}
}

// DeletarPorProdutoID deleta todas as imagens de um produto
func (r *ProductImageRepository) DeletarPorProdutoID(productID int) error {
	query := "DELETE FROM product_images WHERE product_id = ?"
	_, err := r.db.Exec(query, productID)
	if err != nil {
		return fmt.Errorf("erro ao deletar imagens do produto: %w", err)
	}
	return nil
}

// Criar cria uma nova imagem de produto
func (r *ProductImageRepository) Criar(img *models.ProductImage) error {
	query := `
		INSERT INTO product_images (image_type, image_src, product_id, Image_src_small)
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, img.ImageType, img.ImageSrc, img.ProductID, img.ImageSrcSmall)
	if err != nil {
		return fmt.Errorf("erro ao criar imagem: %w", err)
	}

	id, _ := result.LastInsertId()
	img.ID = int(id)
	return nil
}

// ListarPorProdutoID lista todas as imagens de um produto
func (r *ProductImageRepository) ListarPorProdutoID(productID int) ([]models.ProductImage, error) {
	query := `
		SELECT id, image_type, image_src, product_id, Image_src_small
		FROM product_images
		WHERE product_id = ?
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar imagens: %w", err)
	}
	defer rows.Close()

	imagens := []models.ProductImage{}
	for rows.Next() {
		var img models.ProductImage
		err := rows.Scan(&img.ID, &img.ImageType, &img.ImageSrc, &img.ProductID, &img.ImageSrcSmall)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler imagem: %w", err)
		}
		imagens = append(imagens, img)
	}

	return imagens, nil
}
