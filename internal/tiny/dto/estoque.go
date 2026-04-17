package dto

// EstoqueRequest representa a requisição para buscar estoque do produto
type EstoqueRequest struct {
	Token        string `json:"token"`
	ID           string `json:"id"`
	IdListaPreco string `json:"idListaPreco"`
}

// NovoEstoqueRequest cria uma nova requisição para buscar estoque
func NovoEstoqueRequest(token, id, idListaPreco string) *EstoqueRequest {
	return &EstoqueRequest{
		Token:        token,
		ID:           id,
		IdListaPreco: idListaPreco,
	}
}

// EstoqueResponse representa a resposta da API de estoque
type EstoqueResponse struct {
	Retorno RetornoEstoque `json:"retorno"`
}

// RetornoEstoque representa o retorno da API
type RetornoEstoque struct {
	StatusProcessamento string          `json:"status_processamento"`
	Status              string          `json:"status"`
	CodigoErro          string          `json:"codigo_erro,omitempty"`
	Produto             ProdutoEstoque  `json:"produto,omitempty"`
	Erros               []ErroDetalhado `json:"erros,omitempty"`
}

// ProdutoEstoque contém informações de estoque do produto
type ProdutoEstoque struct {
	ID             string            `json:"id"`
	Nome           string            `json:"nome"`
	Codigo         string            `json:"codigo"`
	Unidade        string            `json:"unidade"`
	Saldo          float64           `json:"saldo"`
	SaldoReservado float64           `json:"saldoReservado"`
	Depositos      []DepositoWrapper `json:"depositos,omitempty"`
}

// DepositoWrapper encapsula um depósito
type DepositoWrapper struct {
	Deposito Deposito `json:"deposito"`
}

// Deposito representa um depósito de estoque
type Deposito struct {
	Nome          string  `json:"nome"`
	Desconsiderar string  `json:"desconsiderar"`
	Saldo         float64 `json:"saldo"`
	Empresa       string  `json:"empresa"`
}

// EstoqueTiny é o objeto simplificado para armazenar dados essenciais de estoque
type EstoqueTiny struct {
	IDProduto       string  `json:"id_produto"`
	Nome            string  `json:"nome"`
	Codigo          string  `json:"codigo"`
	Saldo           float64 `json:"saldo"`
	SaldoReservado  float64 `json:"saldo_reservado"`
	SaldoDisponivel float64 `json:"saldo_disponivel"`
}

// ParaEstoqueTiny converte ProdutoEstoque para EstoqueTiny
func (p *ProdutoEstoque) ParaEstoqueTiny() EstoqueTiny {
	return EstoqueTiny{
		IDProduto:       p.ID,
		Nome:            p.Nome,
		Codigo:          p.Codigo,
		Saldo:           p.Saldo,
		SaldoReservado:  p.SaldoReservado,
		SaldoDisponivel: p.Saldo,
	}
}
