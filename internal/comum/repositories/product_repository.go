package repositories

import (
	"database/sql"
	"fmt"

	"prosync/internal/comum/models"
)

// ProductRepository gerencia operações com produtos
type ProductRepository struct {
	db *sql.DB
}

// NovoProductRepository cria novo repositório de produtos
func NovoProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// BuscarPorID busca um produto pelo ID
func (r *ProductRepository) BuscarPorID(id int) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, cost_price, category_id, isEnabled, 
		       isPreSale, sale_count, review_count, sku, stock, observation,
		       weight, height, width, length, product_tiny, ncm, ean, marca, 
		       cest, stop_stock, promotion_id, original_price, created_at, updated_at
		FROM products 
		WHERE id = ?
	`

	var p models.Product
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.CostPrice, &p.CategoryID,
		&p.IsEnabled, &p.IsPreSale, &p.SaleCount, &p.ReviewCount, &p.SKU,
		&p.Stock, &p.Observation, &p.Weight, &p.Height, &p.Width, &p.Length,
		&p.ProductTiny, &p.NCM, &p.EAN, &p.Marca, &p.CEST, &p.StopStock,
		&p.PromotionID, &p.OriginalPrice, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("produto não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar produto: %w", err)
	}

	return &p, nil
}

// BuscarPorSKU busca um produto pelo SKU
func (r *ProductRepository) BuscarPorSKU(sku string) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, cost_price, category_id, isEnabled, 
		       isPreSale, sale_count, review_count, sku, stock, observation,
		       weight, height, width, length, product_tiny, ncm, ean, marca, 
		       cest, stop_stock, promotion_id, original_price, created_at, updated_at
		FROM products 
		WHERE sku = ?
	`

	var p models.Product
	err := r.db.QueryRow(query, sku).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.CostPrice, &p.CategoryID,
		&p.IsEnabled, &p.IsPreSale, &p.SaleCount, &p.ReviewCount, &p.SKU,
		&p.Stock, &p.Observation, &p.Weight, &p.Height, &p.Width, &p.Length,
		&p.ProductTiny, &p.NCM, &p.EAN, &p.Marca, &p.CEST, &p.StopStock,
		&p.PromotionID, &p.OriginalPrice, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Retorna nil sem erro quando não encontra
		}
		return nil, fmt.Errorf("erro ao buscar produto: %w", err)
	}

	return &p, nil
}

// BuscarPorProductTiny busca um produto pelo ID do Tiny
func (r *ProductRepository) BuscarPorProductTiny(tinyID string) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, cost_price, category_id, isEnabled, 
		       isPreSale, sale_count, review_count, sku, stock, observation,
		       weight, height, width, length, product_tiny, ncm, ean, marca, 
		       cest, stop_stock, promotion_id, original_price, created_at, updated_at
		FROM products 
		WHERE product_tiny = ?
	`

	var p models.Product
	err := r.db.QueryRow(query, tinyID).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.CostPrice, &p.CategoryID,
		&p.IsEnabled, &p.IsPreSale, &p.SaleCount, &p.ReviewCount, &p.SKU,
		&p.Stock, &p.Observation, &p.Weight, &p.Height, &p.Width, &p.Length,
		&p.ProductTiny, &p.NCM, &p.EAN, &p.Marca, &p.CEST, &p.StopStock,
		&p.PromotionID, &p.OriginalPrice, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("produto não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar produto: %w", err)
	}

	return &p, nil
}

// ListarTodos lista todos os produtos
func (r *ProductRepository) ListarTodos(limite int) ([]models.Product, error) {
	query := `
		SELECT id, name, description, price, cost_price, category_id, isEnabled, 
		       isPreSale, sale_count, review_count, sku, stock, observation,
		       weight, height, width, length, product_tiny, ncm, ean, marca, 
		       cest, stop_stock, promotion_id, original_price, created_at, updated_at
		FROM products 
		ORDER BY id DESC
	`

	if limite > 0 {
		query += fmt.Sprintf(" LIMIT %d", limite)
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar produtos: %w", err)
	}
	defer rows.Close()

	produtos := []models.Product{}
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Price, &p.CostPrice, &p.CategoryID,
			&p.IsEnabled, &p.IsPreSale, &p.SaleCount, &p.ReviewCount, &p.SKU,
			&p.Stock, &p.Observation, &p.Weight, &p.Height, &p.Width, &p.Length,
			&p.ProductTiny, &p.NCM, &p.EAN, &p.Marca, &p.CEST, &p.StopStock,
			&p.PromotionID, &p.OriginalPrice, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler produto: %w", err)
		}
		produtos = append(produtos, p)
	}

	return produtos, nil
}

// AtualizarEstoque atualiza o estoque de um produto
func (r *ProductRepository) AtualizarEstoque(id int, estoque int) error {
	query := "UPDATE products SET stock = ? WHERE id = ?"
	_, err := r.db.Exec(query, estoque, id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar estoque: %w", err)
	}
	return nil
}

// AtualizarPreco atualiza o preço de um produto
func (r *ProductRepository) AtualizarPreco(id int, preco float64) error {
	query := "UPDATE products SET price = ? WHERE id = ?"
	_, err := r.db.Exec(query, preco, id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar preço: %w", err)
	}
	return nil
}

// Salvar insere ou atualiza um produto
func (r *ProductRepository) Salvar(p *models.Product) error {
	if p.ID == 0 {
		// INSERT
		query := `
			INSERT INTO products (name, description, price, cost_price, category_id, 
			                     isEnabled, isPreSale, sale_count, review_count, sku, 
			                     stock, observation, weight, height, width, length, 
			                     product_tiny, ncm, ean, marca, cest, stop_stock, 
			                     promotion_id, original_price, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		`
		result, err := r.db.Exec(query,
			p.Name, p.Description, p.Price, p.CostPrice, p.CategoryID,
			p.IsEnabled, p.IsPreSale, p.SaleCount, p.ReviewCount, p.SKU,
			p.Stock, p.Observation, p.Weight, p.Height, p.Width, p.Length,
			p.ProductTiny, p.NCM, p.EAN, p.Marca, p.CEST, p.StopStock,
			p.PromotionID, p.OriginalPrice,
		)
		if err != nil {
			return fmt.Errorf("erro ao inserir produto: %w", err)
		}

		id, _ := result.LastInsertId()
		p.ID = int(id)
	} else {
		// UPDATE
		query := `
			UPDATE products SET 
				name = ?, description = ?, price = ?, cost_price = ?, category_id = ?,
				isEnabled = ?, isPreSale = ?, sale_count = ?, review_count = ?, sku = ?,
				stock = ?, observation = ?, weight = ?, height = ?, width = ?, length = ?,
				product_tiny = ?, ncm = ?, ean = ?, marca = ?, cest = ?, stop_stock = ?,
				promotion_id = ?, original_price = ?, updated_at = NOW()
			WHERE id = ?
		`
		_, err := r.db.Exec(query,
			p.Name, p.Description, p.Price, p.CostPrice, p.CategoryID,
			p.IsEnabled, p.IsPreSale, p.SaleCount, p.ReviewCount, p.SKU,
			p.Stock, p.Observation, p.Weight, p.Height, p.Width, p.Length,
			p.ProductTiny, p.NCM, p.EAN, p.Marca, p.CEST, p.StopStock,
			p.PromotionID, p.OriginalPrice, p.ID,
		)
		if err != nil {
			return fmt.Errorf("erro ao atualizar produto: %w", err)
		}
	}

	return nil
}

// CriarOuAtualizar busca produto pelo SKU e atualiza, ou cria novo se não existir
func (r *ProductRepository) CriarOuAtualizar(sku string, p *models.Product) (*models.Product, error) {
	// Tenta buscar produto existente pelo SKU
	produtoExistente, err := r.BuscarPorSKU(sku)

	if err == nil && produtoExistente != nil {
		// Produto existe, atualiza com o ID existente
		p.ID = produtoExistente.ID
		// Mantém campos que não devem ser sobrescritos
		p.SaleCount = produtoExistente.SaleCount
		p.ReviewCount = produtoExistente.ReviewCount
		p.CreatedAt = produtoExistente.CreatedAt
	} else {
		// Produto não existe, será criado (ID = 0)
		p.ID = 0
	}

	// Salva (insert ou update)
	if err := r.Salvar(p); err != nil {
		return nil, err
	}

	return p, nil
}
