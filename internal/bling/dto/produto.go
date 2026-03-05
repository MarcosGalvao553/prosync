package dto

// ProdutoBling representa o payload para criar/atualizar produto no Bling
type ProdutoBling struct {
	Nome                       string      `json:"nome"`
	Codigo                     string      `json:"codigo"`
	Preco                      float64     `json:"preco"`
	Tipo                       string      `json:"tipo"`
	Situacao                   string      `json:"situacao"`
	Formato                    string      `json:"formato"`
	DescricaoCurta             string      `json:"descricaoCurta,omitempty"`
	Unidade                    string      `json:"unidade"`
	PesoLiquido                float64     `json:"pesoLiquido,omitempty"`
	PesoBruto                  float64     `json:"pesoBruto,omitempty"`
	Volumes                    int         `json:"volumes"`
	GTIN                       string      `json:"gtin,omitempty"`
	GTINEmbalagem              string      `json:"gtinEmbalagem,omitempty"`
	Condicao                   int         `json:"condicao"`
	Marca                      string      `json:"marca,omitempty"`
	DescricaoComplementar      string      `json:"descricaoComplementar,omitempty"`
	Observacoes                string      `json:"observacoes,omitempty"`
	DescricaoEmbalagemDiscreta string      `json:"descricaoEmbalagemDiscreta,omitempty"`
	Dimensoes                  *Dimensoes  `json:"dimensoes,omitempty"`
	Tributacao                 *Tributacao `json:"tributacao,omitempty"`
	Midia                      *Midia      `json:"midia,omitempty"`
}

type Dimensoes struct {
	Largura       float64 `json:"largura"`
	Altura        float64 `json:"altura"`
	Profundidade  float64 `json:"profundidade"`
	UnidadeMedida int     `json:"unidadeMedida"` // 1 = Milímetros
}

type Tributacao struct {
	Origem                     int     `json:"origem"`
	NFCI                       string  `json:"nFCI,omitempty"`
	NCM                        string  `json:"ncm,omitempty"`
	CEST                       string  `json:"cest,omitempty"`
	CodigoListaServicos        string  `json:"codigoListaServicos,omitempty"`
	SpedTipoItem               string  `json:"spedTipoItem,omitempty"`
	CodigoItem                 string  `json:"codigoItem,omitempty"`
	PercentualTributos         float64 `json:"percentualTributos"`
	ValorBaseStRetencao        float64 `json:"valorBaseStRetencao"`
	ValorStRetencao            float64 `json:"valorStRetencao"`
	ValorICMSSubstituto        float64 `json:"valorICMSSubstituto"`
	CodigoExcecaoTipi          string  `json:"codigoExcecaoTipi,omitempty"`
	ClasseEnquadramentoIpi     string  `json:"classeEnquadramentoIpi,omitempty"`
	ValorIpiFixo               float64 `json:"valorIpiFixo"`
	CodigoSeloIpi              string  `json:"codigoSeloIpi,omitempty"`
	ValorPisFixo               float64 `json:"valorPisFixo"`
	ValorCofinsFixo            float64 `json:"valorCofinsFixo"`
	CodigoANP                  string  `json:"codigoANP,omitempty"`
	DescricaoANP               string  `json:"descricaoANP,omitempty"`
	PercentualGLP              float64 `json:"percentualGLP"`
	PercentualGasNacional      float64 `json:"percentualGasNacional"`
	PercentualGasImportado     float64 `json:"percentualGasImportado"`
	ValorPartida               float64 `json:"valorPartida"`
	TipoArmamento              int     `json:"tipoArmamento"`
	DescricaoCompletaArmamento string  `json:"descricaoCompletaArmamento,omitempty"`
	DadosAdicionais            string  `json:"dadosAdicionais,omitempty"`
}

type Midia struct {
	Video   *Video   `json:"video,omitempty"`
	Imagens *Imagens `json:"imagens,omitempty"`
}

type Video struct {
	URL string `json:"url,omitempty"`
}

type Imagens struct {
	ImagensURL []ImagemURL `json:"imagensURL,omitempty"`
}

type ImagemURL struct {
	Link string `json:"link"`
}

// RespostaProdutoBling representa a resposta da API do Bling ao buscar/criar produto
type RespostaProdutoBling struct {
	Data ProdutoBlingData `json:"data"`
}

type ProdutoBlingData struct {
	ID       int64   `json:"id"`
	Nome     string  `json:"nome"`
	Codigo   string  `json:"codigo"`
	Preco    float64 `json:"preco"`
	Tipo     string  `json:"tipo"`
	Situacao string  `json:"situacao"`
	// Outros campos conforme necessário
}

// ListaProdutosBling representa a resposta da busca de produtos (com filtro por código)
type ListaProdutosBling struct {
	Data []ProdutoBlingData `json:"data"`
}
