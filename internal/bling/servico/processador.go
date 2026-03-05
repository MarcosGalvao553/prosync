package servico

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	blingDTO "prosync/internal/bling/dto"
	blingEntidade "prosync/internal/bling/entidade"
	"prosync/internal/comum/logger"
	"prosync/internal/comum/models"
	"prosync/internal/comum/repositories"
)

// ItemFilaRateLimit representa um item na fila de retry por rate limit
type ItemFilaRateLimit struct {
	Produto      *models.Product
	SKU          string
	Tentativa    int
	AdicionadoEm time.Time
}

// ProcessadorBling orquestra a sincronização de produtos com o Bling
type ProcessadorBling struct {
	db               *sql.DB
	blingConfigRepo  *repositories.BlingConfigurationRepository
	productUserRepo  *repositories.ProductUserRepository
	productImageRepo *repositories.ProductImageRepository
	logger           *logger.Logger
	filaRateLimit    []ItemFilaRateLimit
	mutexFila        sync.Mutex
}

// NovoProcessadorBling cria uma nova instância do processador
func NovoProcessadorBling(
	db *sql.DB,
	blingConfigRepo *repositories.BlingConfigurationRepository,
	productUserRepo *repositories.ProductUserRepository,
	productImageRepo *repositories.ProductImageRepository,
	logger *logger.Logger,
) *ProcessadorBling {
	return &ProcessadorBling{
		db:               db,
		blingConfigRepo:  blingConfigRepo,
		productUserRepo:  productUserRepo,
		logger:           logger,
		productImageRepo: productImageRepo,
		filaRateLimit:    []ItemFilaRateLimit{},
		mutexFila:        sync.Mutex{},
	}
}

// SincronizarProduto sincroniza um produto com o Bling para todos os usuários vinculados
func (p *ProcessadorBling) SincronizarProduto(produto *models.Product, sku string) error {
	fmt.Printf("\n🔄 Iniciando sincronização Bling para produto ID %d (SKU: %s)\n", produto.ID, sku)

	// Validação: Processar somente produtos ativos
	if !produto.IsEnabled {
		fmt.Printf("   ℹ️  Produto %d está inativo - pulando sincronização Bling\n", produto.ID)
		p.logger.RegistrarInfo("bling",
			fmt.Sprintf("Produto %d (SKU: %s) está inativo - sincronização ignorada", produto.ID, sku),
		)
		return nil
	}

	// Busca todos os usuários vinculados a este produto
	productUsers, err := p.productUserRepo.ListarPorProductID(produto.ID)
	if err != nil {
		return fmt.Errorf("erro ao buscar product_users: %w", err)
	}

	if len(productUsers) == 0 {
		fmt.Printf("   ℹ️  Nenhum usuário vinculado ao produto %d - pulando Bling\n", produto.ID)
		return nil
	}

	fmt.Printf("   👥 %d usuário(s) vinculado(s) ao produto\n", len(productUsers))

	// Processa cada usuário
	for _, pu := range productUsers {
		if err := p.processarProductUser(produto, sku, &pu); err != nil {
			// Verifica se é erro de rate limit
			if isRateLimitError(err) {
				fmt.Printf("   ⏱️  Rate limit detectado - produto adicionado à fila de retry\n")
				p.adicionarNaFilaRateLimit(produto, sku)
				return nil // Não retorna erro, será reprocessado
			}

			p.logger.RegistrarErro("bling",
				fmt.Sprintf("Erro ao processar Bling para user_id %d, product_id %d", pu.UserID, produto.ID),
				err,
			)
			fmt.Printf("   ❌ Erro no Bling para usuário %d: %v\n", pu.UserID, err)
			// Continua para os próximos usuários mesmo com erro
			continue
		}
	}

	return nil
}

// processarProductUser processa a sincronização para um usuário específico
func (p *ProcessadorBling) processarProductUser(produto *models.Product, sku string, pu *models.ProductUser) error {
	// Busca configuração do Bling para este usuário
	config, err := p.blingConfigRepo.BuscarPorUserID(pu.UserID)
	if err != nil {
		return fmt.Errorf("erro ao buscar config Bling: %w", err)
	}

	if config == nil {
		fmt.Printf("   ⚠️  Usuário %d sem configuração Bling - pulando\n", pu.UserID)
		return nil
	}

	// Cria cliente Bling
	client := blingEntidade.NovoBlingClient(
		config.ClientID.String,
		config.SecretKey.String,
		config.AccessToken.String,
		p.logger,
	)

	// Valida e atualiza token se necessário
	if err := p.validarEAtualizarToken(client, config); err != nil {
		return fmt.Errorf("erro ao validar token: %w", err)
	}

	// Busca imagens do produto
	imagens, err := p.productImageRepo.ListarPorProdutoID(produto.ID)
	if err != nil {
		p.logger.RegistrarErro("bling",
			fmt.Sprintf("Erro ao buscar imagens do produto %d", produto.ID),
			err,
		)
		imagens = []models.ProductImage{} // Continua sem imagens
	}

	// Busca produto no Bling pelo código/SKU
	produtoBling, err := p.buscarProdutoBling(client, sku, pu)
	if err != nil {
		return fmt.Errorf("erro ao buscar produto no Bling: %w", err)
	}

	var blingProductID int64

	if produtoBling != nil {
		// Produto existe - atualizar
		blingProductID = produtoBling.ID
		if err := p.atualizarProdutoBling(client, produtoBling.ID, produto, imagens, pu, sku); err != nil {
			return fmt.Errorf("erro ao atualizar produto no Bling: %w", err)
		}

		// Atualiza bling_product_id na product_user se ainda não estiver preenchido
		if !pu.BlingProductID.Valid || pu.BlingProductID.String == "" {
			if err := p.productUserRepo.AtualizarBlingProductID(pu.ID, fmt.Sprintf("%d", produtoBling.ID)); err != nil {
				p.logger.RegistrarErro("bling",
					fmt.Sprintf("Erro ao atualizar bling_product_id na product_user %d", pu.ID),
					err,
				)
			}
		}

		fmt.Printf("   ✅ Produto atualizado no Bling (user %d, bling_id %d)\n", pu.UserID, produtoBling.ID)
	} else {
		// Produto não existe - criar
		novoProduto, err := p.criarProdutoBling(client, produto, imagens, pu, sku)
		if err != nil {
			return fmt.Errorf("erro ao criar produto no Bling: %w", err)
		}
		blingProductID = novoProduto.ID

		// Atualiza product_user com o bling_product_id
		if err := p.productUserRepo.AtualizarBlingProductID(pu.ID, fmt.Sprintf("%d", blingProductID)); err != nil {
			p.logger.RegistrarErro("bling",
				fmt.Sprintf("Erro ao atualizar bling_product_id na product_user %d", pu.ID),
				err,
			)
		}

		fmt.Printf("   ✅ Produto criado no Bling (user %d, bling_id %d)\n", pu.UserID, blingProductID)
	}

	// Atualiza estoque no Bling
	if err := p.atualizarEstoqueBling(client, blingProductID, produto, pu); err != nil {
		return fmt.Errorf("erro ao atualizar estoque no Bling: %w", err)
	}

	return nil
}

// validarEAtualizarToken verifica se o token está válido e renova se necessário
func (p *ProcessadorBling) validarEAtualizarToken(client *blingEntidade.BlingClient, config *models.BlingConfiguration) error {
	if p.blingConfigRepo.TokenEstaValido(config) {
		fmt.Printf("   ✅ Token do Bling válido para usuário %d\n", config.UserID.Int64)
		p.logger.RegistrarInfo("bling", fmt.Sprintf("Token válido para user_id %d", config.UserID.Int64))
		return nil
	}

	// Token inválido ou perto de expirar - faz refresh
	fmt.Printf("   🔄 Renovando token do Bling para usuário %d...\n", config.UserID.Int64)

	tokenResp, err := client.RefreshToken(config.RefreshToken.String)
	if err != nil {
		p.logger.RegistrarErro("bling",
			fmt.Sprintf("Erro ao renovar token para user_id %d", config.UserID.Int64),
			err,
		)

		// Registra log específico de atualização de token
		p.logger.RegistrarChamada(logger.EntradaLog{
			Servico:  "bling",
			Operacao: "AtualizaTokenBling",
			UserID:   uint64(config.UserID.Int64),
			Requisicao: map[string]interface{}{
				"user_id":       config.UserID.Int64,
				"grant_type":    "refresh_token",
				"refresh_token": "***", // Ocultado por segurança
			},
			Erro: err.Error(),
		})

		return err
	}

	// Calcula nova data de expiração
	novaExpiracao := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	// Atualiza tokens no banco
	if err := p.blingConfigRepo.AtualizarTokens(
		config.ID,
		tokenResp.AccessToken,
		tokenResp.RefreshToken,
		novaExpiracao,
	); err != nil {
		return fmt.Errorf("erro ao salvar novos tokens: %w", err)
	}

	// Atualiza o client com o novo token
	client.SetAccessToken(tokenResp.AccessToken)

	// Registra log de sucesso
	p.logger.RegistrarChamada(logger.EntradaLog{
		Servico:  "bling",
		Operacao: "AtualizaTokenBling",
		UserID:   uint64(config.UserID.Int64),
		Requisicao: map[string]interface{}{
			"user_id":       config.UserID.Int64,
			"grant_type":    "refresh_token",
			"refresh_token": "***", // Ocultado por segurança
		},
		Resposta: map[string]interface{}{
			"expires_in": tokenResp.ExpiresIn,
			"user_id":    config.UserID.Int64,
		},
	})

	fmt.Printf("   ✅ Token renovado com sucesso\n")
	return nil
}

// buscarProdutoBling busca um produto no Bling pelo código/SKU
func (p *ProcessadorBling) buscarProdutoBling(client *blingEntidade.BlingClient, codigo string, pu *models.ProductUser) (*blingDTO.ProdutoBlingData, error) {
	produto, err := client.BuscarProdutoPorCodigo(codigo)

	// Registra log da busca
	logEntry := logger.EntradaLog{
		Servico:  "bling",
		Operacao: "BuscaProdutoBling",
		SKU:      codigo,
		UserID:   uint64(pu.UserID),
		Requisicao: map[string]interface{}{
			"codigo": codigo,
		},
	}

	if err != nil {
		logEntry.Erro = err.Error()
		p.logger.RegistrarChamada(logEntry)
		return nil, err
	}

	if produto != nil {
		logEntry.Resposta = map[string]interface{}{
			"id":         produto.ID,
			"nome":       produto.Nome,
			"codigo":     produto.Codigo,
			"encontrado": true,
		}
	} else {
		logEntry.Resposta = map[string]interface{}{
			"encontrado": false,
		}
	}
	p.logger.RegistrarChamada(logEntry)

	return produto, nil
}

// criarProdutoBling cria um novo produto no Bling
func (p *ProcessadorBling) criarProdutoBling(
	client *blingEntidade.BlingClient,
	produto *models.Product,
	imagens []models.ProductImage,
	pu *models.ProductUser,
	sku string,
) (*blingDTO.ProdutoBlingData, error) {

	payload := p.formatarProdutoParaCriar(produto, imagens, pu)

	novoProduto, err := client.CriarProduto(payload)

	// Registra log
	logEntry := logger.EntradaLog{
		Servico:  "bling",
		Operacao: "CriaProdutoBling",
		SKU:      sku,
		UserID:   uint64(pu.UserID),
		Requisicao: map[string]interface{}{
			"nome":             payload.Nome,
			"codigo":           payload.Codigo,
			"preco":            payload.Preco,
			"payload_completo": payload,
		},
	}

	if err != nil {
		logEntry.Erro = err.Error()
		p.logger.RegistrarChamada(logEntry)
		return nil, err
	}

	logEntry.Resposta = map[string]interface{}{
		"id":     novoProduto.ID,
		"nome":   novoProduto.Nome,
		"codigo": novoProduto.Codigo,
	}
	p.logger.RegistrarChamada(logEntry)

	return novoProduto, nil
}

// atualizarProdutoBling atualiza um produto existente no Bling
func (p *ProcessadorBling) atualizarProdutoBling(
	client *blingEntidade.BlingClient,
	blingProductID int64,
	produto *models.Product,
	imagens []models.ProductImage,
	pu *models.ProductUser,
	sku string,
) error {

	payload := p.formatarProdutoParaAtualizar(produto, imagens, pu)

	err := client.AtualizarProduto(blingProductID, payload)

	// Registra log
	logEntry := logger.EntradaLog{
		Servico:  "bling",
		Operacao: "AtualizaProdutoBling",
		SKU:      sku,
		UserID:   uint64(pu.UserID),
		Requisicao: map[string]interface{}{
			"bling_product_id": blingProductID,
			"nome":             payload.Nome,
			"codigo":           payload.Codigo,
			"preco":            payload.Preco,
			"payload_completo": payload,
		},
	}

	if err != nil {
		logEntry.Erro = err.Error()
		p.logger.RegistrarChamada(logEntry)
		return err
	}

	logEntry.Resposta = map[string]interface{}{
		"bling_product_id": blingProductID,
	}
	p.logger.RegistrarChamada(logEntry)

	return nil
}

// atualizarEstoqueBling atualiza o estoque de um produto no Bling
func (p *ProcessadorBling) atualizarEstoqueBling(
	client *blingEntidade.BlingClient,
	blingProductID int64,
	produto *models.Product,
	pu *models.ProductUser,
) error {

	// Busca deposito padrão
	depositos, err := client.BuscarDepositos()
	if err != nil {
		return fmt.Errorf("erro ao buscar depósitos: %w", err)
	}

	var deposito *blingDTO.DepositoBling
	for i := range depositos {
		if depositos[i].Padrao {
			deposito = &depositos[i]
			break
		}
	}

	if deposito == nil {
		return fmt.Errorf("depósito padrão não encontrado")
	}

	// Prepara payload de estoque
	estoque := &blingDTO.EstoqueBling{
		Produto: blingDTO.ProdutoEstoque{
			ID: blingProductID,
		},
		Deposito: blingDTO.DepositoEstoque{
			ID: deposito.ID,
		},
		Operacao:    "B", // Balanço
		Quantidade:  float64(produto.Stock.Int64),
		Preco:       0.0,
		Custo:       0.0,
		Observacoes: "Atualizado NerdDrop",
	}

	err = client.AtualizarEstoque(estoque)

	// Registra log
	logEntry := logger.EntradaLog{
		Servico:  "bling",
		Operacao: "AtualizaEstoqueBling",
		SKU:      produto.SKU.String,
		UserID:   uint64(pu.UserID),
		Requisicao: map[string]interface{}{
			"bling_product_id": blingProductID,
			"deposito_id":      deposito.ID,
			"quantidade":       produto.Stock.Int64,
			"operacao":         "B",
			"payload_completo": estoque,
		},
	}

	if err != nil {
		logEntry.Erro = err.Error()
		p.logger.RegistrarChamada(logEntry)
		return err
	}

	logEntry.Resposta = map[string]interface{}{
		"bling_product_id": blingProductID,
		"quantidade":       produto.Stock.Int64,
		"deposito_id":      deposito.ID,
	}
	p.logger.RegistrarChamada(logEntry)

	fmt.Printf("   📦 Estoque atualizado no Bling: %d unidades\n", produto.Stock.Int64)
	return nil
}

// formatarProdutoParaCriar formata um produto para criação no Bling
func (p *ProcessadorBling) formatarProdutoParaCriar(
	produto *models.Product,
	imagens []models.ProductImage,
	pu *models.ProductUser,
) *blingDTO.ProdutoBling {

	// Prepara imagens
	imagensURL := []blingDTO.ImagemURL{}
	for _, img := range imagens {
		imagensURL = append(imagensURL, blingDTO.ImagemURL{
			Link: img.ImageSrc,
		})
	}

	// Determina preço (usa price do product_user se disponível, senão usa do produto)
	preco := produto.Price.Float64
	if pu.Price.Valid {
		preco = pu.Price.Float64
	}

	return &blingDTO.ProdutoBling{
		Nome:                       produto.Name,
		Codigo:                     produto.SKU.String,
		Preco:                      preco,
		Tipo:                       "P", // Produto
		Situacao:                   "A", // Ativo
		Formato:                    "S", // Simples
		DescricaoCurta:             produto.Description.String,
		Unidade:                    "UN",
		PesoLiquido:                produto.Weight.Float64,
		PesoBruto:                  produto.Weight.Float64,
		Volumes:                    1,
		GTIN:                       produto.EAN.String,
		GTINEmbalagem:              produto.EAN.String,
		Condicao:                   0,
		Marca:                      produto.Marca.String,
		DescricaoComplementar:      produto.Description.String,
		Observacoes:                "",
		DescricaoEmbalagemDiscreta: "",
		Dimensoes: &blingDTO.Dimensoes{
			Largura:       produto.Width.Float64,
			Altura:        produto.Height.Float64,
			Profundidade:  produto.Length.Float64,
			UnidadeMedida: 1, // Milímetros
		},
		Tributacao: &blingDTO.Tributacao{
			Origem:                 0,
			NCM:                    produto.NCM.String,
			CEST:                   produto.CEST.String,
			PercentualTributos:     0,
			ValorBaseStRetencao:    0,
			ValorStRetencao:        0,
			ValorICMSSubstituto:    0,
			ValorIpiFixo:           0,
			ValorPisFixo:           0,
			ValorCofinsFixo:        0,
			PercentualGLP:          0,
			PercentualGasNacional:  0,
			PercentualGasImportado: 0,
			ValorPartida:           0,
			TipoArmamento:          0,
		},
		Midia: &blingDTO.Midia{
			Video: &blingDTO.Video{
				URL: "",
			},
			Imagens: &blingDTO.Imagens{
				ImagensURL: imagensURL,
			},
		},
	}
}

// formatarProdutoParaAtualizar formata um produto para atualização no Bling
func (p *ProcessadorBling) formatarProdutoParaAtualizar(
	produto *models.Product,
	imagens []models.ProductImage,
	pu *models.ProductUser,
) *blingDTO.ProdutoBling {

	// Prepara imagens
	imagensURL := []blingDTO.ImagemURL{}
	for _, img := range imagens {
		imagensURL = append(imagensURL, blingDTO.ImagemURL{
			Link: img.ImageSrc,
		})
	}

	// Determina preço
	preco := produto.Price.Float64
	if pu.Price.Valid {
		preco = pu.Price.Float64
	}

	return &blingDTO.ProdutoBling{
		Nome:                  produto.Name,
		Codigo:                produto.SKU.String,
		Preco:                 preco,
		Tipo:                  "P",
		Situacao:              "A",
		Formato:               "S",
		DescricaoCurta:        produto.Description.String,
		Unidade:               "UN",
		PesoLiquido:           produto.Weight.Float64,
		PesoBruto:             produto.Weight.Float64,
		Volumes:               1,
		GTIN:                  produto.EAN.String,
		GTINEmbalagem:         produto.EAN.String,
		Condicao:              0,
		Marca:                 produto.Marca.String,
		DescricaoComplementar: produto.Description.String,
		Dimensoes: &blingDTO.Dimensoes{
			Largura:       produto.Width.Float64,
			Altura:        produto.Height.Float64,
			Profundidade:  produto.Length.Float64,
			UnidadeMedida: 1,
		},
		Tributacao: &blingDTO.Tributacao{
			Origem: 0,
			NCM:    produto.NCM.String,
			CEST:   produto.CEST.String,
		},
		Midia: &blingDTO.Midia{
			Video: &blingDTO.Video{
				URL: "",
			},
			Imagens: &blingDTO.Imagens{
				ImagensURL: imagensURL,
			},
		},
	}
}

// adicionarNaFilaRateLimit adiciona um produto na fila de retry por rate limit
func (p *ProcessadorBling) adicionarNaFilaRateLimit(produto *models.Product, sku string) {
	p.mutexFila.Lock()
	defer p.mutexFila.Unlock()

	// Verifica se já está na fila
	for _, item := range p.filaRateLimit {
		if item.Produto.ID == produto.ID && item.SKU == sku {
			return // Já está na fila
		}
	}

	p.filaRateLimit = append(p.filaRateLimit, ItemFilaRateLimit{
		Produto:      produto,
		SKU:          sku,
		Tentativa:    0,
		AdicionadoEm: time.Now(),
	})

	p.logger.RegistrarInfo("bling",
		fmt.Sprintf("Produto %d (SKU: %s) adicionado à fila de retry por rate limit. Total na fila: %d",
			produto.ID, sku, len(p.filaRateLimit)))
}

// ProcessarFilaRateLimit processa produtos que falharam por rate limit
func (p *ProcessadorBling) ProcessarFilaRateLimit() {
	p.mutexFila.Lock()
	itens := make([]ItemFilaRateLimit, len(p.filaRateLimit))
	copy(itens, p.filaRateLimit)
	p.filaRateLimit = []ItemFilaRateLimit{} // Limpa a fila
	p.mutexFila.Unlock()

	if len(itens) == 0 {
		return
	}

	fmt.Printf("\n🔄 Processando fila de rate limit: %d produto(s)\n", len(itens))

	for _, item := range itens {
		// Aguarda 1 minuto desde que foi adicionado
		tempoDecorrido := time.Since(item.AdicionadoEm)
		if tempoDecorrido < time.Minute {
			aguardar := time.Minute - tempoDecorrido
			fmt.Printf("   ⏳ Aguardando %v antes de reprocessar produto %d\n", aguardar.Round(time.Second), item.Produto.ID)
			time.Sleep(aguardar)
		}

		fmt.Printf("   🔁 Retry %d - Produto %d (SKU: %s)\n", item.Tentativa+1, item.Produto.ID, item.SKU)

		err := p.SincronizarProduto(item.Produto, item.SKU)
		if err != nil {
			if isRateLimitError(err) && item.Tentativa < 3 {
				// Se ainda é rate limit e não excedeu tentativas, adiciona novamente
				p.mutexFila.Lock()
				item.Tentativa++
				item.AdicionadoEm = time.Now()
				p.filaRateLimit = append(p.filaRateLimit, item)
				p.mutexFila.Unlock()

				fmt.Printf("   ⚠️  Rate limit novamente - tentativa %d/3\n", item.Tentativa)
			} else {
				p.logger.RegistrarErro("bling",
					fmt.Sprintf("Erro após retry do produto %d (tentativa %d)", item.Produto.ID, item.Tentativa+1),
					err,
				)
				fmt.Printf("   ❌ Falha definitiva após %d tentativas\n", item.Tentativa+1)
			}
		} else {
			fmt.Printf("   ✅ Produto %d reprocessado com sucesso\n", item.Produto.ID)
		}
	}
}

// TemItensNaFilaRateLimit verifica se há itens na fila
func (p *ProcessadorBling) TemItensNaFilaRateLimit() bool {
	p.mutexFila.Lock()
	defer p.mutexFila.Unlock()
	return len(p.filaRateLimit) > 0
}

// isRateLimitError verifica se o erro é de rate limit
func isRateLimitError(err error) bool {
	if err == nil {
		return false
	}

	// Verifica se é o erro customizado RateLimitError
	if _, ok := err.(*blingEntidade.RateLimitError); ok {
		return true
	}

	// Verifica se a mensagem contém "rate limit"
	errMsg := err.Error()
	return contains(errMsg, "Rate limit") || contains(errMsg, "rate limit") || contains(errMsg, "429")
}

// contains verifica se uma string contém outra (case-insensitive helper)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
