package servico

import (
	"database/sql"
	"fmt"
	"math"

	"prosync/internal/comum/logger"
	"prosync/internal/comum/models"
	"prosync/internal/comum/repositories"
	"prosync/internal/trovata/dto"
	"prosync/internal/trovata/entidade"
)

// ProcessadorTrovata gerencia a sincronização com Trovata
type ProcessadorTrovata struct {
	client      *entidade.TrovataClient
	logger      *logger.Logger
	partnerRepo *repositories.PartnerRepository
	partnerID   int
}

// NovoProcessadorTrovata cria nova instância do processador
func NovoProcessadorTrovata(
	db *sql.DB,
	logger *logger.Logger,
) *ProcessadorTrovata {
	return &ProcessadorTrovata{
		client:      entidade.NovoTrovataClient(logger),
		logger:      logger,
		partnerRepo: repositories.NovoPartnerRepository(db),
		partnerID:   1, // Por enquanto fixo em 1
	}
}

// SincronizarProduto sincroniza produto com Trovata (criar produto + atualizar estoque)
func (p *ProcessadorTrovata) SincronizarProduto(produto *models.Product, categoria *models.Category, sku string, idProdutoTiny string) error {
	p.logger.RegistrarInfo("trovata",
		fmt.Sprintf("Iniciando sincronização do produto ID %d (SKU: %s, Tiny ID: %s) com Trovata", produto.ID, sku, idProdutoTiny),
	)

	// Busca partner para cálculo de preço
	partner, err := p.partnerRepo.BuscarPorID(p.partnerID)
	if err != nil {
		p.logger.RegistrarErro("trovata",
			fmt.Sprintf("Erro ao buscar partner %d", p.partnerID),
			err,
		)
		return fmt.Errorf("erro ao buscar partner: %w", err)
	}

	// 1. Criar produto na Trovata
	if err := p.criarProduto(produto, categoria, partner, sku, idProdutoTiny); err != nil {
		return fmt.Errorf("erro ao criar produto na Trovata: %w", err)
	}

	// 2. Atualizar estoque do produto
	if err := p.atualizarEstoque(produto, sku, idProdutoTiny); err != nil {
		return fmt.Errorf("erro ao atualizar estoque na Trovata: %w", err)
	}

	p.logger.RegistrarInfo("trovata",
		fmt.Sprintf("Produto ID %d sincronizado com sucesso na Trovata", produto.ID),
	)

	return nil
}

// criarProduto cria o produto na Trovata
func (p *ProcessadorTrovata) criarProduto(produto *models.Product, categoria *models.Category, partner *models.Partner, sku string, idProdutoTiny string) error {
	// Calcula preço usando a mesma lógica do PHP
	preco := p.calcularPreco(produto.Price.Float64, partner)

	// Monta descrição (limita a 249 caracteres)
	descricao2 := ""
	if produto.Description.Valid {
		desc := produto.Description.String
		if len(desc) > 249 {
			descricao2 = desc[:249]
		} else {
			descricao2 = desc
		}
	}

	// Nome da categoria
	nomeCategoria := ""
	if categoria != nil {
		nomeCategoria = categoria.Name
	}

	// Monta o request
	request := &dto.ProdutoTrovataRequest{
		Produto:                 produto.ID,
		DescricaoProduto:        produto.Name,
		ApelidoProduto:          obterString(produto.SKU),
		AbreviaturaUnidade:      nil,
		GrupoProduto:            nil,
		SubgrupoProduto:         nil,
		Situacao:                obterSituacao(produto.IsEnabled),
		PesoLiquido:             nil,
		ClassificacaoFiscal:     obterString(produto.NCM),
		Categoria:               nomeCategoria,
		PontoCritico:            nil,
		Grade:                   nil,
		CodigoBarras:            nil,
		Especificacao:           nomeCategoria,
		PrecoBase:               preco,
		FamiliaComercial:        nil,
		UnidadeFabricacao:       nil,
		Especie:                 nil,
		Segmento:                nil,
		TipoEmbalagem:           nil,
		UsoProdutoOpcional:      nil,
		DescricaoProduto2:       descricao2,
		DescricaoProduto3:       nil,
		Marca:                   obterString(produto.Marca),
		TipoProduto:             nil,
		EstiloUso:               nil,
		DimensaoTamanho:         nil,
		Nicho:                   nil,
		Linha:                   nil,
		Genero:                  nil,
		NCM:                     obterString(produto.NCM),
		PrecoCusto:              preco,
		PrecoFinal:              preco,
		ListaMultiploVenda:      nil,
		GradePor:                nil,
		SubstituicaoTributaria:  nil,
		PercDesconto:            nil,
		PercDescontoParceria:    nil,
		PercDescontoGerencial:   nil,
		PercDescontoPromocional: nil,
		Colecao:                 nil,
		ValidaEstoque:           nil,
	}

	return p.client.CriarProduto(request, sku, idProdutoTiny)
}

// atualizarEstoque atualiza o estoque do produto na Trovata
func (p *ProcessadorTrovata) atualizarEstoque(produto *models.Product, sku string, idProdutoTiny string) error {
	estoque := 0
	if produto.Stock.Valid {
		estoque = int(produto.Stock.Int64)
	}

	request := &dto.EstoqueTrovataRequest{
		SaldoEstoque:          produto.ID,
		Produto:               produto.ID,
		Complemento1:          nil,
		Complemento2:          nil,
		LocalEstoque:          nil,
		Tipo:                  "1",
		Ano:                   nil,
		Mes:                   nil,
		SaldoInicial:          nil,
		Entradas:              nil,
		Saidas:                nil,
		Complemento3:          nil,
		SaldoFinal:            estoque,
		SaldoDisponivel:       estoque,
		SaldoAuxiliar:         estoque,
		ReservaERP:            "0.000",
		ReservaOnline:         "0.000",
		ReservaLocal:          "0.000",
		PontoCritico:          "0.000",
		DataDisponivelInicial: nil,
		DataDisponivelFinal:   nil,
		PeriodoEntregaInicial: nil,
		PeriodoEntregaFinal:   nil,
		Dia:                   nil,
		DataBaseSaldoEstoque:  nil,
	}

	return p.client.AtualizarEstoque(request, sku, idProdutoTiny)
}

// calcularPreco calcula o preço para Trovata usando a mesma lógica do PHP
func (p *ProcessadorTrovata) calcularPreco(preco float64, partner *models.Partner) float64 {
	// valor com desconto de 5%
	vp := preco - (preco * 0.05)

	// taxa = (100 - fee_adicional - desconto_negociado) / 100
	taxa := (100 - partner.FeeAdicional - partner.DescontoNegociado) / 100

	// final = vp / taxa
	final := vp / taxa

	// Arredonda para 2 casas decimais
	return math.Round(final*100) / 100
}

// obterString converte sql.NullString para string
func obterString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// obterSituacao retorna "A" (ativo) ou "I" (inativo)
func obterSituacao(isEnabled bool) string {
	if isEnabled {
		return "A"
	}
	return "I"
}
