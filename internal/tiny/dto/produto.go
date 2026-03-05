package dto

// ProdutoRequest representa a requisição para buscar dados do produto
type ProdutoRequest struct {
	Token   string `json:"token"`
	ID      string `json:"id"`
	Formato string `json:"formato"`
}

// NovoProdutoRequest cria uma nova requisição para buscar produto
func NovoProdutoRequest(token, id string) *ProdutoRequest {
	return &ProdutoRequest{
		Token:   token,
		ID:      id,
		Formato: "json",
	}
}

// ProdutoResponse representa a resposta da API de produto
type ProdutoResponse struct {
	Retorno RetornoProduto `json:"retorno"`
}

// RetornoProduto representa o retorno da API
type RetornoProduto struct {
	StatusProcessamento string          `json:"status_processamento"`
	Status              string          `json:"status"`
	CodigoErro          string          `json:"codigo_erro,omitempty"`
	Produto             DadosProduto    `json:"produto,omitempty"`
	Erros               []ErroDetalhado `json:"erros,omitempty"`
}

// ErroDetalhado representa um erro retornado pela API
type ErroDetalhado struct {
	Erro string `json:"erro"`
}

// DadosProduto contém as informações do produto
type DadosProduto struct {
	ID                    string   `json:"id"`
	Nome                  string   `json:"nome"`
	Codigo                string   `json:"codigo"`
	Unidade               string   `json:"unidade"`
	Preco                 float64  `json:"preco"`
	PrecoPromocional      float64  `json:"preco_promocional"`
	NCM                   string   `json:"ncm"`
	Origem                string   `json:"origem"`
	GTIN                  string   `json:"gtin"`
	GTINEmbalagem         string   `json:"gtin_embalagem"`
	Localizacao           string   `json:"localizacao"`
	PesoLiquido           float64  `json:"peso_liquido"`
	PesoBruto             float64  `json:"peso_bruto"`
	EstoqueMinimo         float64  `json:"estoque_minimo"`
	EstoqueMaximo         float64  `json:"estoque_maximo"`
	IDFornecedor          int      `json:"id_fornecedor"`
	NomeFornecedor        string   `json:"nome_fornecedor"`
	CodigoFornecedor      string   `json:"codigo_fornecedor"`
	CodigoPeloFornecedor  string   `json:"codigo_pelo_fornecedor"`
	UnidadePorCaixa       string   `json:"unidade_por_caixa"`
	PrecoCusto            float64  `json:"preco_custo"`
	PrecoCustoMedio       float64  `json:"preco_custo_medio"`
	Situacao              string   `json:"situacao"`
	Tipo                  string   `json:"tipo"`
	ClasseIPI             string   `json:"classe_ipi"`
	ValorIPIFixo          string   `json:"valor_ipi_fixo"`
	CodListaServicos      string   `json:"cod_lista_servicos"`
	DescricaoComplementar string   `json:"descricao_complementar"`
	Garantia              string   `json:"garantia"`
	CEST                  string   `json:"cest"`
	Obs                   string   `json:"obs"`
	TipoVariacao          string   `json:"tipoVariacao"`
	Variacoes             string   `json:"variacoes"`
	IDProdutoPai          string   `json:"idProdutoPai"`
	SobEncomenda          string   `json:"sob_encomenda"`
	DiasPreparacao        string   `json:"dias_preparacao"`
	Marca                 string   `json:"marca"`
	TipoEmbalagem         string   `json:"tipoEmbalagem"`
	AlturaEmbalagem       string   `json:"alturaEmbalagem"`
	ComprimentoEmbalagem  string   `json:"comprimentoEmbalagem"`
	LarguraEmbalagem      string   `json:"larguraEmbalagem"`
	DiametroEmbalagem     string   `json:"diametroEmbalagem"`
	QtdVolumes            string   `json:"qtd_volumes"`
	Categoria             string   `json:"categoria"`
	Anexos                []Anexo  `json:"anexos,omitempty"`
	ImagensExternas       []string `json:"imagens_externas,omitempty"`
	ClasseProduto         string   `json:"classe_produto"`
	SEOTitle              string   `json:"seo_title"`
	SEOKeywords           string   `json:"seo_keywords"`
	LinkVideo             string   `json:"link_video"`
	SEODescription        string   `json:"seo_description"`
	Slug                  string   `json:"slug"`
}

// Anexo representa um anexo do produto
type Anexo struct {
	URL string `json:"anexo"`
}

// ProdutoTiny é o objeto simplificado para armazenar dados essenciais do produto
type ProdutoTiny struct {
	ID                    string  `json:"id"`
	Nome                  string  `json:"nome"`
	Codigo                string  `json:"codigo"` // Código do produto
	Preco                 float64 `json:"preco"`
	PrecoCusto            float64 `json:"preco_custo"`
	Marca                 string  `json:"marca"`
	Categoria             string  `json:"categoria"`
	Situacao              string  `json:"situacao"`
	GTIN                  string  `json:"gtin"`
	NCM                   string  `json:"ncm"`
	CEST                  string  `json:"cest"`
	DescricaoComplementar string  `json:"descricao_complementar"`
	Obs                   string  `json:"obs"`
	PesoLiquido           float64 `json:"peso_liquido"`
	AlturaEmbalagem       string  `json:"altura_embalagem"`
	ComprimentoEmbalagem  string  `json:"comprimento_embalagem"`
	LarguraEmbalagem      string  `json:"largura_embalagem"`
	Anexos                []Anexo `json:"anexos,omitempty"`
}

// ParaProdutoTiny converte DadosProduto para ProdutoTiny
func (d *DadosProduto) ParaProdutoTiny() ProdutoTiny {
	return ProdutoTiny{
		ID:                    d.ID,
		Nome:                  d.Nome,
		Codigo:                d.Codigo,
		Preco:                 d.Preco,
		PrecoCusto:            d.PrecoCusto,
		Marca:                 d.Marca,
		Categoria:             d.Categoria,
		Situacao:              d.Situacao,
		GTIN:                  d.GTIN,
		NCM:                   d.NCM,
		CEST:                  d.CEST,
		DescricaoComplementar: d.DescricaoComplementar,
		Obs:                   d.Obs,
		PesoLiquido:           d.PesoLiquido,
		AlturaEmbalagem:       d.AlturaEmbalagem,
		ComprimentoEmbalagem:  d.ComprimentoEmbalagem,
		LarguraEmbalagem:      d.LarguraEmbalagem,
		Anexos:                d.Anexos,
	}
}
