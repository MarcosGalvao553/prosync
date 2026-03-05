package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"prosync/internal/comum/database"
	"prosync/internal/comum/models"
)

// ExemploUso demonstra como usar os modelos e repositórios
// Esta função pode ser chamada de outros pacotes para testar
func ExemploUso() {
	// Conecta ao banco de dados
	db := database.ObterConexao()

	// Cria o repositório de produtos
	produtoRepo := NovoProductRepository(db)

	// Exemplo 1: Buscar produto por ID
	produto, err := produtoRepo.BuscarPorID(1)
	if err != nil {
		log.Printf("Erro ao buscar produto: %v", err)
	} else {
		fmt.Printf("Produto encontrado: %s (ID: %d)\n", produto.Name, produto.ID)

		// Acessa campos nullable
		if produto.SKU.Valid {
			fmt.Printf("SKU: %s\n", produto.SKU.String)
		}

		if produto.Price.Valid {
			fmt.Printf("Preço: R$ %.2f\n", produto.Price.Float64)
		}

		if produto.Stock.Valid {
			fmt.Printf("Estoque: %d\n", produto.Stock.Int64)
		}
	}

	// Exemplo 2: Buscar produto pelo Product Tiny ID
	produtoPorTiny, err := produtoRepo.BuscarPorProductTiny("123456")
	if err != nil {
		log.Printf("Erro ao buscar produto por Tiny ID: %v", err)
	} else {
		fmt.Printf("Produto Tiny: %s\n", produtoPorTiny.Name)
	}

	// Exemplo 3: Listar produtos (limite de 10)
	produtos, err := produtoRepo.ListarTodos(10)
	if err != nil {
		log.Printf("Erro ao listar produtos: %v", err)
	} else {
		fmt.Printf("\nTotal de produtos: %d\n", len(produtos))
		for i, p := range produtos {
			fmt.Printf("%d. %s (ID: %d)\n", i+1, p.Name, p.ID)
		}
	}

	// Exemplo 4: Atualizar estoque
	err = produtoRepo.AtualizarEstoque(1, 100)
	if err != nil {
		log.Printf("Erro ao atualizar estoque: %v", err)
	} else {
		fmt.Println("Estoque atualizado com sucesso!")
	}

	// Exemplo 5: Atualizar preço
	err = produtoRepo.AtualizarPreco(1, 99.90)
	if err != nil {
		log.Printf("Erro ao atualizar preço: %v", err)
	} else {
		fmt.Println("Preço atualizado com sucesso!")
	}

	// Exemplo 6: Criar novo produto
	novoProduto := &models.Product{
		Name:        "Produto Teste",
		CostPrice:   50.00,
		CategoryID:  1,
		IsEnabled:   true,
		SaleCount:   sql.NullInt64{Int64: 0, Valid: true},
		ReviewCount: sql.NullInt64{Int64: 0, Valid: true},
	}

	err = produtoRepo.Salvar(novoProduto)
	if err != nil {
		log.Printf("Erro ao salvar produto: %v", err)
	} else {
		fmt.Printf("Produto criado com ID: %d\n", novoProduto.ID)
	}

	// Exemplo 7: Atualizar produto existente
	novoProduto.Name = "Produto Teste Atualizado"
	err = produtoRepo.Salvar(novoProduto)
	if err != nil {
		log.Printf("Erro ao atualizar produto: %v", err)
	} else {
		fmt.Println("Produto atualizado com sucesso!")
	}
}
