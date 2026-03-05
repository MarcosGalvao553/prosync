package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"prosync/internal/comum/models"
)

// CategoryRepository gerencia operações com categorias
type CategoryRepository struct {
	db *sql.DB
}

// NovoCategoryRepository cria novo repositório de categorias
func NovoCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// BuscarPorNome busca uma categoria pelo nome
func (r *CategoryRepository) BuscarPorNome(nome string) (*models.Category, error) {
	query := `
		SELECT id, name, image_src, code, category_id, range_1, range_2, range_3, free_shipping
		FROM categories 
		WHERE name = ?
		LIMIT 1
	`

	var c models.Category
	err := r.db.QueryRow(query, nome).Scan(
		&c.ID, &c.Name, &c.ImageSrc, &c.Code, &c.CategoryID,
		&c.Range1, &c.Range2, &c.Range3, &c.FreeShipping,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Retorna nil se não encontrar
		}
		return nil, fmt.Errorf("erro ao buscar categoria: %w", err)
	}

	return &c, nil
}

// Criar cria uma nova categoria
func (r *CategoryRepository) Criar(nome string) (*models.Category, error) {
	query := "INSERT INTO categories (name) VALUES (?)"

	result, err := r.db.Exec(query, nome)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar categoria: %w", err)
	}

	id, _ := result.LastInsertId()

	return &models.Category{
		ID:   int(id),
		Name: nome,
	}, nil
}

// BuscarOuCriarPorNome busca uma categoria pelo nome ou cria se não existir
func (r *CategoryRepository) BuscarOuCriarPorNome(nome string) (*models.Category, error) {
	// Remove espaços em branco extras
	nome = strings.TrimSpace(nome)

	if nome == "" {
		return nil, fmt.Errorf("nome da categoria não pode ser vazio")
	}

	// Tenta buscar primeiro
	categoria, err := r.BuscarPorNome(nome)
	if err != nil {
		return nil, err
	}

	// Se encontrou, retorna
	if categoria != nil {
		return categoria, nil
	}

	// Se não encontrou, cria
	return r.Criar(nome)
}

// ProcessarCategoriaTiny processa a categoria do Tiny no formato "Cat1>>Cat2"
// Retorna a subcategoria (Cat2) se existir, caso contrário a categoria principal (Cat1)
func (r *CategoryRepository) ProcessarCategoriaTiny(categoriaTiny string) (*models.Category, error) {
	if categoriaTiny == "" {
		return nil, fmt.Errorf("categoria do Tiny está vazia")
	}

	// Divide por ">>"
	partes := strings.Split(categoriaTiny, ">>")

	var nomeCategoria string
	if len(partes) > 1 {
		// Se tem subcategoria, usa ela
		nomeCategoria = strings.TrimSpace(partes[1])
	} else {
		// Senão, usa a primeira
		nomeCategoria = strings.TrimSpace(partes[0])
	}

	// Busca ou cria a categoria
	return r.BuscarOuCriarPorNome(nomeCategoria)
}

// ListarTodas lista todas as categorias
func (r *CategoryRepository) ListarTodas() ([]models.Category, error) {
	query := `
		SELECT id, name, image_src, code, category_id, range_1, range_2, range_3, free_shipping
		FROM categories 
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar categorias: %w", err)
	}
	defer rows.Close()

	categorias := []models.Category{}
	for rows.Next() {
		var c models.Category
		err := rows.Scan(
			&c.ID, &c.Name, &c.ImageSrc, &c.Code, &c.CategoryID,
			&c.Range1, &c.Range2, &c.Range3, &c.FreeShipping,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler categoria: %w", err)
		}
		categorias = append(categorias, c)
	}

	return categorias, nil
}
