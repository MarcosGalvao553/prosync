package main

import (
	"fmt"
	"log"
	"time"

	blingServico "prosync/internal/bling/servico"
	"prosync/internal/comum/config"
	"prosync/internal/comum/database"
	"prosync/internal/comum/logger"
	"prosync/internal/comum/repositories"
	"prosync/internal/comum/servidor"
	"prosync/internal/tiny/entidade"
	"prosync/internal/tiny/servico"
)

func main() {
	fmt.Println("=== ProSync - Iniciando ===")

	// Carrega configurações
	cfg, err := config.CarregarConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}
	fmt.Printf("Configurações carregadas. Ambiente: %s\n", cfg.Ambiente)

	// Inicializa conexão com banco de dados
	if err := database.InicializarConexao(cfg); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	fmt.Println("✅ Conectado ao banco de dados")
	defer database.FecharConexao()

	// Inicializa o logger
	logger, err := logger.NovoLogger()
	if err != nil {
		log.Fatalf("Erro ao inicializar logger: %v", err)
	}
	fmt.Println("Logger inicializado")

	// Cria o cliente Tiny
	tinyClient := entidade.NovoTinyClient(cfg, logger)
	fmt.Println("Cliente Tiny criado")

	// Cria os repositórios
	db := database.ObterConexao()
	categoryRepo := repositories.NovoCategoryRepository(db)
	productRepo := repositories.NovoProductRepository(db)
	productPromotionRepo := repositories.NovoProductPromotionRepository(db)
	productImageRepo := repositories.NovoProductImageRepository(db)
	blingConfigRepo := repositories.NovoBlingConfigurationRepository(db)
	productUserRepo := repositories.NovoProductUserRepository(db)
	fmt.Println("Repositórios criados")

	// Cria o processador Bling
	processadorBling := blingServico.NovoProcessadorBling(
		db,
		blingConfigRepo,
		productUserRepo,
		productImageRepo,
		logger,
	)

	// Cria o processador
	processador := servico.NovoProcessadorTiny(tinyClient, logger, categoryRepo, productRepo, productPromotionRepo, productImageRepo, processadorBling)
	fmt.Println("Processador criado")

	// Inicia servidor web em goroutine
	servidorWeb := servidor.NovoServidorWeb("8080", logger)
	go func() {
		if err := servidorWeb.Iniciar(); err != nil {
			log.Printf("Erro no servidor web: %v", err)
		}
	}()

	// Executa o processamento em loop
	fmt.Printf("Intervalo de execução: %d minutos\n", cfg.IntervaloExecucaoMinutos)

	for {
		fmt.Printf("\n[%s] Iniciando ciclo de processamento...\n", time.Now().Format("2006-01-02 15:04:05"))

		if err := processar(processador, processadorBling, logger); err != nil {
			log.Printf("Erro no processamento: %v", err)
			logger.RegistrarErro("sistema", "Erro no ciclo de processamento", err)
		}

		// Aguarda o intervalo configurado antes da próxima execução
		proximaExecucao := time.Now().Add(time.Duration(cfg.IntervaloExecucaoMinutos) * time.Minute)
		fmt.Printf("[%s] Próxima execução em: %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			proximaExecucao.Format("2006-01-02 15:04:05"))

		time.Sleep(time.Duration(cfg.IntervaloExecucaoMinutos) * time.Minute)
	}
}

// processar executa o ciclo completo de processamento
func processar(processador *servico.ProcessadorTiny, processadorBling *blingServico.ProcessadorBling, logger *logger.Logger) error {
	logger.RegistrarInfo("sistema", "=== INÍCIO DO CICLO DE PROCESSAMENTO ===")

	// Processa todas as exceções de lista de preços
	produtos, err := processador.ProcessarExcecoesListaPreco()

	if err != nil {
		return fmt.Errorf("erro ao processar exceções: %w", err)
	}

	// Obtém estatísticas do processamento
	stats := processador.EstatisticasProcessamento(produtos)

	// Exibe estatísticas
	fmt.Println("\n=== ESTATÍSTICAS DO PROCESSAMENTO ===")
	fmt.Printf("Total de produtos: %d\n", stats["total"])
	fmt.Printf("Com dados completos: %d\n", stats["com_dados"])
	fmt.Printf("Com estoque: %d\n", stats["com_estoque"])
	fmt.Printf("Totalmente completos: %d\n", stats["completos"])
	fmt.Printf("Com estoque disponível: %d\n", stats["com_estoque_positivo"])
	fmt.Printf("Valor total em estoque: R$ %.2f\n", stats["valor_total_estoque"])

	// Log das estatísticas
	logger.RegistrarInfo("sistema", fmt.Sprintf(
		"Processamento concluído - Total: %d | Completos: %d | Com estoque: %d | Valor: R$ %.2f",
		stats["total"], stats["completos"], stats["com_estoque_positivo"], stats["valor_total_estoque"],
	))

	// Exemplo: mostra os primeiros produtos completos
	exibirExemplosProdutos(produtos, 3)

	// Processa fila de rate limit do Bling APÓS processar todos os produtos
	// Isso evita concorrência e garante que o rate limit (1 req/s) seja respeitado
	if processadorBling.TemItensNaFilaRateLimit() {
		fmt.Println("\n🔄 === PROCESSANDO FILA DE RATE LIMIT ===")
		logger.RegistrarInfo("sistema", "Iniciando processamento da fila de rate limit")
		processadorBling.ProcessarFilaRateLimit()
		fmt.Println("✅ Fila de rate limit processada")
	}

	logger.RegistrarInfo("sistema", "=== FIM DO CICLO DE PROCESSAMENTO ===")
	return nil
}

// exibirExemplosProdutos mostra alguns exemplos de produtos processados
func exibirExemplosProdutos(produtos []servico.ProdutoCompleto, limite int) {
	fmt.Println("\n=== EXEMPLOS DE PRODUTOS ===")

	contador := 0
	for _, prod := range produtos {
		// Mostra apenas produtos completos
		if prod.Produto != nil && prod.Estoque != nil {
			fmt.Printf("\nProduto: %s\n", prod.Produto.Nome)
			fmt.Printf("  ID: %s | Código: %s\n", prod.Produto.ID, prod.Produto.Codigo)
			fmt.Printf("  Preço Lista: R$ %.2f | Preço Exceção: R$ %.2f\n",
				prod.Produto.Preco, prod.Excecao.Preco)
			fmt.Printf("  Estoque: %.0f | Reservado: %.0f | Disponível: %.0f\n",
				prod.Estoque.Saldo, prod.Estoque.SaldoReservado, prod.Estoque.SaldoDisponivel)
			fmt.Printf("  Marca: %s | Categoria: %s\n", prod.Produto.Marca, prod.Produto.Categoria)

			contador++
			if contador >= limite {
				break
			}
		}
	}

	if contador == 0 {
		fmt.Println("Nenhum produto completo para exibir")
	}
}
