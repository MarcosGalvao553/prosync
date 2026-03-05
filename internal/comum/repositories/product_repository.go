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
	fmt.Printf("   🔎 [DEBUG] BuscarPorSKU - Query com SKU: '%s'\n", sku)

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
			fmt.Printf("   ℹ️  [DEBUG] Nenhum registro encontrado para SKU: '%s'\n", sku)
			return nil, nil // Retorna nil sem erro quando não encontra
		}
		fmt.Printf("   ❌ [DEBUG] Erro na query: %v\n", err)
		return nil, fmt.Errorf("erro ao buscar produto: %w", err)
	}

	skuValue := "NULL"
	if p.SKU.Valid {
		skuValue = p.SKU.String
	}
	fmt.Printf("   ✅ [DEBUG] Produto encontrado - ID: %d, SKU no banco: '%s'\n", p.ID, skuValue)

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
	fmt.Printf("\n🔍 [DEBUG] CriarOuAtualizar - Buscando produto com SKU: %s\n", sku)

	// Tenta buscar produto existente pelo SKU
	produtoExistente, err := r.BuscarPorSKU(sku)

	if err != nil {
		fmt.Printf("❌ [DEBUG] Erro ao buscar produto: %v\n", err)
	}

	if err == nil && produtoExistente != nil {
		// Produto existe, atualiza com o ID existente
		fmt.Printf("✅ [DEBUG] Produto ENCONTRADO - ID: %d, SKU: %s\n", produtoExistente.ID, sku)
		fmt.Printf("   → Será ATUALIZADO (não criará duplicata)\n")

		p.ID = produtoExistente.ID
		// Mantém campos que não devem ser sobrescritos
		p.SaleCount = produtoExistente.SaleCount
		p.ReviewCount = produtoExistente.ReviewCount
		p.CreatedAt = produtoExistente.CreatedAt
	} else {
		// Produto não existe, será criado (ID = 0)
		fmt.Printf("⚠️  [DEBUG] Produto NÃO encontrado - SKU: %s\n", sku)
		fmt.Printf("   → Será CRIADO novo registro\n")
		p.ID = 0
	}

	fmt.Printf("💾 [DEBUG] Chamando Salvar com ID: %d\n", p.ID)

	// Salva (insert ou update)
	if err := r.Salvar(p); err != nil {
		fmt.Printf("❌ [DEBUG] Erro ao salvar: %v\n", err)
		return nil, err
	}

	fmt.Printf("✅ [DEBUG] Produto salvo com sucesso - ID final: %d\n", p.ID)
	return p, nil
}
