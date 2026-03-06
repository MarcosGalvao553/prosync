package dto

// ProdutoTrovataRequest representa a requisição de criação de produto na Trovata
type ProdutoTrovataRequest struct {
	Produto              int     `json:"PRODUTO"`
	DescricaoProduto     string  `json:"DESCRICAO_PRODUTO"`
	ApelidoProduto       string  `json:"APELIDO_PRODUTO"`
	AbreviaturaUnidade   *string `json:"ABREVIATURA_UNIDADE"`
	GrupoProduto         *string `json:"GRUPO_PRODUTO"`
	SubgrupoProduto      *string `json:"SUBGRUPO_PRODUTO"`
	Situacao             string  `json:"SITUACAO"`
	PesoLiquido          *string `json:"PESO_LIQUIDO"`
	ClassificacaoFiscal  string  `json:"CLASSIFICACAO_FISCAL"`
	Categoria            string  `json:"CATEGORIA"`
	PontoCritico         *string `json:"PONTO_CRITICO"`
	Grade                *string `json:"GRADE"`
	CodigoBarras         *string `json:"CODIGO_BARRAS"`
	Especificacao        string  `json:"ESPECIFICACAO"`
	PrecoBase            float64 `json:"PRECO_BASE"`
	FamiliaComercial     *string `json:"FAMILIA_COMERCIAL"`
	UnidadeFabricacao    *string `json:"UNIDADE_FABRICACAO"`
	Especie              *string `json:"ESPECIE"`
	Segmento             *string `json:"SEGMENTO"`
	TipoEmbalagem        *string `json:"TIPO_EMBALAGEM"`
	UsoProdutoOpcional   *string `json:"USO_PRODUTO_OPCIONAL"`
	DescricaoProduto2    string  `json:"DESCRICAO_PRODUTO_2"`
	DescricaoProduto3    *string `json:"DESCRICAO_PRODUTO_3"`
	Marca                string  `json:"MARCA"`
	TipoProduto          *string `json:"TIPO_PRODUTO"`
	EstiloUso            *string `json:"ESTILO_USO"`
	DimensaoTamanho      *string `json:"DIMENSAO_TAMANHO"`
	Nicho                *string `json:"NICHO"`
	Linha                *string `json:"LINHA"`
	Genero               *string `json:"GENERO"`
	NCM                  string  `json:"NCM"`
	PrecoCusto           float64 `json:"PRECO_CUSTO"`
	PrecoFinal           float64 `json:"PRECO_FINAL"`
	ListaMultiploVenda   *string `json:"LISTA_MULTIPLO_VENDA"`
	GradePor             *string `json:"GRADE_POR"`
	SubstituicaoTributaria *string `json:"SUBSTITUICAO_TRIBUTARIA"`
	PercDesconto         *string `json:"PERC_DESCONTO"`
	PercDescontoParceria *string `json:"PERC_DESCONTO_PARCERIA"`
	PercDescontoGerencial *string `json:"PERC_DESCONTO_GERENCIAL"`
	PercDescontoPromocional *string `json:"PERC_DESCONTO_PROMOCIONAL"`
	Colecao              *string `json:"COLECAO"`
	ValidaEstoque        *string `json:"VALIDA_ESTOQUE"`
}

// EstoqueTrovataRequest representa a requisição de atualização de estoque na Trovata
type EstoqueTrovataRequest struct {
	SaldoEstoque         int     `json:"SALDO_ESTOQUE"`
	Produto              int     `json:"PRODUTO"`
	Complemento1         *string `json:"COMPLEMENTO_1"`
	Complemento2         *string `json:"COMPLEMENTO_2"`
	LocalEstoque         *string `json:"LOCAL_ESTOQUE"`
	Tipo                 string  `json:"TIPO"`
	Ano                  *string `json:"ANO"`
	Mes                  *string `json:"MES"`
	SaldoInicial         *string `json:"SALDO_INICIAL"`
	Entradas             *string `json:"ENTRADAS"`
	Saidas               *string `json:"SAIDAS"`
	Complemento3         *string `json:"COMPLEMENTO_3"`
	SaldoFinal           int     `json:"SALDO_FINAL"`
	SaldoDisponivel      int     `json:"SALDO_DISPONIVEL"`
	SaldoAuxiliar        int     `json:"SALDO_AUXILIAR"`
	ReservaERP           string  `json:"RESERVA_ERP"`
	ReservaOnline        string  `json:"RESERVA_ONLINE"`
	ReservaLocal         string  `json:"RESERVA_LOCAL"`
	PontoCritico         string  `json:"PONTO_CRITICO"`
	DataDisponivelInicial *string `json:"DATA_DISPONIVEL_INICIAL"`
	DataDisponivelFinal  *string `json:"DATA_DISPONIVEL_FINAL"`
	PeriodoEntregaInicial *string `json:"PERIODO_ENTREGA_INICIAL"`
	PeriodoEntregaFinal  *string `json:"PERIODO_ENTREGA_FINAL"`
	Dia                  *string `json:"DIA"`
	DataBaseSaldoEstoque *string `json:"DATA_BASE_SALDO_ESTOQUE"`
}

// OcorrenciaVendaRequest representa a requisição de atualização de status do pedido
type OcorrenciaVendaRequest struct {
	Observacao      string `json:"OBSERVACAO"`
	PedidoVendedor  string `json:"PEDIDO_VENDEDOR"`
	SituacaoVenda   string `json:"SITUACAO_VENDA"`
}

// TrovataResponse representa a resposta genérica da API Trovata
type TrovataResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}
