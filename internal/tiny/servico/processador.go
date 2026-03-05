package servico

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	blingServico "prosync/internal/bling/servico"
	"prosync/internal/comum/logger"
	"prosync/internal/comum/models"
	"prosync/internal/comum/repositories"
	"prosync/internal/tiny/dto"
	"prosync/internal/tiny/entidade"
)

// ProdutoCompleto agrega todas as informações de um produto
type ProdutoCompleto struct {
	Excecao          dto.ProdutoExcecaoListaPrecoTiny
	Produto          *dto.ProdutoTiny
	Estoque          *dto.EstoqueTiny
	Categoria        *models.Category
	EstoqueCalculado float64         // Estoque calculado para Nerdrop
	ProdutoNerdrop   *models.Product // Produto no banco Nerdrop (existente ou novo)
	ProdutoParams    *models.Product // Parâmetros preparados para salvar
}

// ProcessadorTiny orquestra o processamento dos dados do Tiny
type ProcessadorTiny struct {
	client               *entidade.TinyClient
	logger               *logger.Logger
	categoryRepo         *repositories.CategoryRepository
	productRepo          *repositories.ProductRepository
	productPromotionRepo *repositories.ProductPromotionRepository
	productImageRepo     *repositories.ProductImageRepository
	processadorBling     *blingServico.ProcessadorBling // Opcional
}

const (
	// DecreaseStock é a quantidade a ser subtraída do estoque disponível
	DecreaseStock = 3
	// MaintenancePrice é o valor fixo de manutenção adicionado ao preço
	MaintenancePrice = 3.50
)

// NovoProcessadorTiny cria uma nova instância do processador
func NovoProcessadorTiny(client *entidade.TinyClient, logger *logger.Logger, categoryRepo *repositories.CategoryRepository, productRepo *repositories.ProductRepository, productPromotionRepo *repositories.ProductPromotionRepository, productImageRepo *repositories.ProductImageRepository, processadorBling *blingServico.ProcessadorBling) *ProcessadorTiny {
	return &ProcessadorTiny{
		client:               client,
		logger:               logger,
		productRepo:          productRepo,
		categoryRepo:         categoryRepo,
		productPromotionRepo: productPromotionRepo,
		productImageRepo:     productImageRepo,
		processadorBling:     processadorBling,
	}
}

// ProcessarExcecoesListaPreco executa o fluxo completo:
// 1. Busca todas as exceções de lista de preços
// 2. Para cada produto, busca dados e estoque
// 3. Retorna uma coleção com todas as informações
func (p *ProcessadorTiny) ProcessarExcecoesListaPreco() ([]ProdutoCompleto, error) {
	p.logger.RegistrarInfo("processador", "=== INÍCIO DO PROCESSAMENTO DE EXCEÇÕES ===")

	// Passo 1: Buscar todas as exceções de lista de preços
	p.logger.RegistrarInfo("processador", "Buscando exceções de lista de preços...")
	excecoes, err := p.client.BuscarTodasExcecoesListaPreco()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar exceções: %w", err)
	}

	p.logger.RegistrarInfo("processador", fmt.Sprintf("Total de %d exceções encontradas", len(excecoes)))

	// Passo 2: Para cada exceção, buscar dados e estoque do produto
	produtosCompletos := make([]ProdutoCompleto, 0, len(excecoes))

	for i, excecao := range excecoes {
		idProduto := strconv.FormatInt(excecao.IdProduto, 10)

		p.logger.RegistrarInfo("processador", fmt.Sprintf(
			"[%d/%d] Processando produto ID %s...",
			i+1, len(excecoes), idProduto,
		))

		produtoCompleto := ProdutoCompleto{
			Excecao: excecao,
		}

		// Busca dados do produto
		produto, errProduto := p.client.BuscarDadosProduto(idProduto)
		if errProduto != nil {
			p.logger.RegistrarErro("processador",
				fmt.Sprintf("Erro ao buscar dados do produto %s", idProduto),
				errProduto,
			)
			// Continua processamento mesmo com erro
			continue
		}

		// Validação: Processar somente produtos ativos (situação "A")
		if produto.Situacao != "A" {
			p.logger.RegistrarInfo("processador",
				fmt.Sprintf("Produto %s - '%s' com situação '%s' (inativo), não será processado",
					idProduto, produto.Nome, produto.Situacao,
				),
			)
			continue
		}

		// Validação: Se produto não for FUNKO ou BLOKEES, não criar
		nomeUpper := strings.ToUpper(produto.Nome)
		isFunko := strings.Contains(nomeUpper, "FUNKO")
		isBlokees := strings.Contains(nomeUpper, "BLOKEES")

		if !isFunko && !isBlokees {
			p.logger.RegistrarInfo("processador",
				fmt.Sprintf("Produto %s - '%s' não é FUNKO/BLOKEES, não será processado",
					idProduto, produto.Nome,
				),
			)
			continue
		}

		produtoCompleto.Produto = produto

		// Processa a categoria do produto
		if produto.Categoria != "" {
			categoria, errCategoria := p.categoryRepo.ProcessarCategoriaTiny(produto.Categoria)
			if errCategoria != nil {
				p.logger.RegistrarErro("processador",
					fmt.Sprintf("Erro ao processar categoria '%s' do produto %s", produto.Categoria, idProduto),
					errCategoria,
				)
			} else {
				produtoCompleto.Categoria = categoria
				p.logger.RegistrarInfo("processador",
					fmt.Sprintf("Categoria processada: %s (ID: %d)", categoria.Name, categoria.ID),
				)
			}
		}

		// Busca estoque do produto
		estoque, errEstoque := p.client.BuscarEstoqueProduto(idProduto)
		if errEstoque != nil {
			p.logger.RegistrarErro("processador",
				fmt.Sprintf("Erro ao buscar estoque do produto %s", idProduto),
				errEstoque,
			)
			// Continua processamento mesmo com erro
		} else {
			produtoCompleto.Estoque = estoque

			// Calcula estoque disponível para Nerdrop
			// Fórmula: saldo - saldo_reservado - decrease_stock (3)
			saldo := estoque.SaldoDisponivel
			saldoReservado := estoque.SaldoReservado
			estoqueCalculado := saldo - saldoReservado - DecreaseStock

			// Se negativo, ajusta para 0
			if estoqueCalculado < 0 {
				estoqueCalculado = 0
			}

			produtoCompleto.EstoqueCalculado = estoqueCalculado

			// Registra log do cálculo de estoque
			p.logger.RegistrarChamada(logger.EntradaLog{
				Servico:       "nerdrop",
				Operacao:      "CalculoEstoqueNerdrop",
				ProdutoTinyID: idProduto,
				SKU:           produto.Codigo,
				Resposta: map[string]interface{}{
					"id_produto_tiny":   idProduto,
					"sku":               produto.Codigo,
					"saldo_tiny":        saldo,
					"saldo_reservado":   saldoReservado,
					"decrease_stock":    DecreaseStock,
					"estoque_calculado": estoqueCalculado,
					"formula":           "saldo - saldo_reservado - decrease_stock",
				},
			})

			p.logger.RegistrarInfo("processador",
				fmt.Sprintf("Produto %s - Estoque calculado: %.0f (Saldo: %.0f, Reservado: %.0f, Decrease: %d)",
					idProduto, estoqueCalculado, saldo, saldoReservado, DecreaseStock,
				),
			)

			// Busca produto no banco pelo SKU
			if produto.Codigo != "" {
				dropProduct, errDrop := p.productRepo.BuscarPorSKU(produto.Codigo)
				if errDrop == nil && dropProduct != nil {
					produtoCompleto.ProdutoNerdrop = dropProduct
					p.logger.RegistrarInfo("processador",
						fmt.Sprintf("Produto encontrado no banco: SKU %s (ID: %d)", produto.Codigo, dropProduct.ID),
					)
				} else {
					p.logger.RegistrarInfo("processador",
						fmt.Sprintf("Produto não encontrado no banco: SKU %s - Será criado novo", produto.Codigo),
					)
				}
			}

			// Calcula o preço do produto
			var precoFinal float64
			var temPromocaoAtiva bool

			// Verifica se o produto tem promoção ativa (somente se o produto existir no banco)
			if produtoCompleto.ProdutoNerdrop != nil && produtoCompleto.ProdutoNerdrop.ID > 0 {
				temPromocao, errPromocao := p.productPromotionRepo.VerificarPromocaoAtiva(produtoCompleto.ProdutoNerdrop.ID)
				if errPromocao == nil {
					temPromocaoAtiva = temPromocao
				}
			}

			// Se não tiver promoção ativa, calcula o novo preço
			if !temPromocaoAtiva {
				// Usa o preço da exceção de lista de preços + valor fixo de manutenção (3.50)
				precoFinal = excecao.Preco + MaintenancePrice

				// Registra log do cálculo de preço
				p.logger.RegistrarChamada(logger.EntradaLog{
					Servico:       "nerdrop",
					Operacao:      "CalculoPrecoNerdrop",
					ProdutoTinyID: idProduto,
					SKU:           produto.Codigo,
					Resposta: map[string]interface{}{
						"id_produto_tiny":  idProduto,
						"sku":              produto.Codigo,
						"preco_tiny":       excecao.Preco,
						"valor_manutencao": MaintenancePrice,
						"preco_calculado":  precoFinal,
						"formula":          "preco_tiny + valor_manutencao",
					},
				})

				p.logger.RegistrarInfo("processador",
					fmt.Sprintf("Produto %s - Preço Tiny: %.2f + Manutenção: %.2f = Total: %.2f",
						idProduto, excecao.Preco, MaintenancePrice, precoFinal),
				)
			} else {
				p.logger.RegistrarInfo("processador",
					fmt.Sprintf("Produto %s - Tem promoção ativa, preço não será atualizado", idProduto),
				)
				// Mantém o preço atual do produto
				if produtoCompleto.ProdutoNerdrop != nil && produtoCompleto.ProdutoNerdrop.Price.Valid {
					precoFinal = produtoCompleto.ProdutoNerdrop.Price.Float64
				}
			}

			// Prepara parâmetros do produto para salvar no banco
			produtoParams := p.prepararParametrosProduto(produto, produtoCompleto.Categoria, estoqueCalculado, precoFinal, idProduto)
			produtoCompleto.ProdutoParams = produtoParams

			// Registra log com os dados do produto preparados
			produtoJSON, _ := json.Marshal(produtoParams)
			var produtoMap map[string]interface{}
			json.Unmarshal(produtoJSON, &produtoMap)

			p.logger.RegistrarChamada(logger.EntradaLog{
				Servico:       "nerdrop",
				Operacao:      "DadosProdutoNerdrop",
				ProdutoTinyID: idProduto,
				SKU:           produto.Codigo,
				Resposta:      produtoMap,
			})

			// Verifica se produto já existe para mostrar mensagem apropriada
			produtoExistente, _ := p.productRepo.BuscarPorSKU(produto.Codigo)
			isNovo := produtoExistente == nil

			// Salva ou atualiza o produto no banco
			if isNovo {
				fmt.Printf("\n💾 Criando novo produto %s (SKU: %s)...\n", idProduto, produto.Codigo)
			} else {
				fmt.Printf("\n🔄 Atualizando produto %s (SKU: %s, ID: %d)...\n", idProduto, produto.Codigo, produtoExistente.ID)
			}

			produtoSalvo, errSalvar := p.productRepo.CriarOuAtualizar(produto.Codigo, produtoParams)
			if errSalvar != nil {
				p.logger.RegistrarErro("processador",
					fmt.Sprintf("Erro ao salvar produto %s (SKU: %s)", idProduto, produto.Codigo),
					errSalvar,
				)
				fmt.Printf("❌ Erro ao salvar produto %s (SKU: %s): %v\n", idProduto, produto.Codigo, errSalvar)
			} else {
				produtoCompleto.ProdutoNerdrop = produtoSalvo
				if isNovo {
					p.logger.RegistrarInfo("processador",
						fmt.Sprintf("Produto criado com sucesso: ID %d, SKU %s", produtoSalvo.ID, produto.Codigo),
					)
					fmt.Printf("✅ Produto criado: ID %d, SKU %s, Nome: %s\n", produtoSalvo.ID, produto.Codigo, produto.Nome)
				} else {
					p.logger.RegistrarInfo("processador",
						fmt.Sprintf("Produto atualizado com sucesso: ID %d, SKU %s", produtoSalvo.ID, produto.Codigo),
					)
					fmt.Printf("✅ Produto atualizado: ID %d, SKU %s, Nome: %s\n", produtoSalvo.ID, produto.Codigo, produto.Nome)
				}

				// Deleta imagens antigas do produto
				if err := p.productImageRepo.DeletarPorProdutoID(produtoSalvo.ID); err != nil {
					p.logger.RegistrarErro("processador",
						fmt.Sprintf("Erro ao deletar imagens antigas do produto %d", produtoSalvo.ID),
						err,
					)
					fmt.Printf("   ⚠️  Erro ao deletar imagens antigas do produto %d\n", produtoSalvo.ID)
				}

				// Salva novas imagens do produto
				if len(produto.Anexos) > 0 {
					for _, anexo := range produto.Anexos {
						img := &models.ProductImage{
							ImageType: 0,
							ImageSrc:  anexo.URL,
							ProductID: produtoSalvo.ID,
						}

						if err := p.productImageRepo.Criar(img); err != nil {
							p.logger.RegistrarErro("processador",
								fmt.Sprintf("Erro ao criar imagem para produto %d: %s", produtoSalvo.ID, anexo.URL),
								err,
							)
						}
					}

					p.logger.RegistrarInfo("processador",
						fmt.Sprintf("Produto %d - %d imagens salvas", produtoSalvo.ID, len(produto.Anexos)),
					)
					fmt.Printf("   📸 %d imagens salvas para produto %d\n", len(produto.Anexos), produtoSalvo.ID)
				}
				// Sincroniza com Bling (se processadorBling estiver configurado)
				if p.processadorBling != nil {
					if err := p.processadorBling.SincronizarProduto(produtoSalvo, produto.Codigo); err != nil {
						p.logger.RegistrarErro("processador",
							fmt.Sprintf("Erro ao sincronizar produto %d com Bling", produtoSalvo.ID),
							err,
						)
						fmt.Printf("   ⚠️  Erro na sincronização Bling: %v\n", err)
					}
				}
			}
		}

		produtosCompletos = append(produtosCompletos, produtoCompleto)
	}

	p.logger.RegistrarInfo("processador", fmt.Sprintf(
		"=== PROCESSAMENTO CONCLUÍDO - %d produtos processados ===",
		len(produtosCompletos),
	))

	return produtosCompletos, nil
}

// prepararParametrosProduto prepara os parâmetros do produto para salvar no banco
func (p *ProcessadorTiny) prepararParametrosProduto(produto *dto.ProdutoTiny, categoria *models.Category, estoqueCalculado float64, preco float64, produtoTinyID string) *models.Product {

	toNullString := func(s string) sql.NullString {
		if s == "" {
			return sql.NullString{Valid: false}
		}
		return sql.NullString{String: s, Valid: true}
	}

	toNullFloat64 := func(s string) sql.NullFloat64 {
		if s == "" {
			return sql.NullFloat64{Valid: false}
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return sql.NullFloat64{Valid: false}
		}
		return sql.NullFloat64{Float64: f, Valid: true}
	}

	toNullInt64 := func(f float64) sql.NullInt64 {
		return sql.NullInt64{Int64: int64(f), Valid: true}
	}

	// Determina se o produto está ativo (situação "A" = ativo)
	isEnabled := produto.Situacao == "A"

	// ID da categoria (se existir)
	categoryID := 74 // Default
	if categoria != nil {
		categoryID = categoria.ID
	}

	// Monta o produto
	productParams := &models.Product{
		Name:        produto.Nome,
		Description: toNullString(produto.DescricaoComplementar),
		Price:       sql.NullFloat64{Float64: preco, Valid: preco > 0},
		CostPrice:   produto.PrecoCusto,
		CategoryID:  categoryID,
		IsEnabled:   isEnabled,
		IsPreSale:   false, // Default false
		SaleCount:   0,
		ReviewCount: 0,
		SKU:         toNullString(produto.Codigo),
		Stock:       toNullInt64(estoqueCalculado),
		Observation: toNullString(produto.Obs),
		Weight:      sql.NullFloat64{Float64: produto.PesoLiquido, Valid: produto.PesoLiquido > 0},
		Height:      toNullFloat64(produto.AlturaEmbalagem),
		Width:       toNullFloat64(produto.ComprimentoEmbalagem),
		Length:      toNullFloat64(produto.LarguraEmbalagem),
		ProductTiny: toNullString(produtoTinyID),
		EAN:         toNullString(produto.GTIN),
		NCM:         toNullString(produto.NCM),
		Marca:       toNullString(produto.Marca),
		CEST:        toNullString(produto.CEST),
		StopStock:   sql.NullInt64{Int64: 0, Valid: true},
	}

	return productParams
}

// EstatisticasProcessamento retorna estatísticas sobre os produtos processados
func (p *ProcessadorTiny) EstatisticasProcessamento(produtos []ProdutoCompleto) map[string]interface{} {
	stats := map[string]interface{}{
		"total":                len(produtos),
		"com_dados":            0,
		"com_estoque":          0,
		"completos":            0,
		"com_estoque_positivo": 0,
		"valor_total_estoque":  0.0,
	}

	for _, prod := range produtos {
		if prod.Produto != nil {
			stats["com_dados"] = stats["com_dados"].(int) + 1
		}
		if prod.Estoque != nil {
			stats["com_estoque"] = stats["com_estoque"].(int) + 1

			if prod.Estoque.SaldoDisponivel > 0 {
				stats["com_estoque_positivo"] = stats["com_estoque_positivo"].(int) + 1
			}

			// Calcula valor total em estoque (saldo disponível * preço da exceção)
			valorProduto := prod.Estoque.SaldoDisponivel * prod.Excecao.Preco
			stats["valor_total_estoque"] = stats["valor_total_estoque"].(float64) + valorProduto
		}
		if prod.Produto != nil && prod.Estoque != nil {
			stats["completos"] = stats["completos"].(int) + 1
		}
	}

	return stats
}
