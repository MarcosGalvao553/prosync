package dto

// EstoqueBling representa o payload para atualizar estoque no Bling
type EstoqueBling struct {
	Produto     ProdutoEstoque  `json:"produto"`
	Deposito    DepositoEstoque `json:"deposito"`
	Operacao    string          `json:"operacao"` // B = Balanço
	Quantidade  float64         `json:"quantidade"`
	Preco       float64         `json:"preco"`
	Custo       float64         `json:"custo"`
	Observacoes string          `json:"observacoes"`
}

type ProdutoEstoque struct {
	ID int64 `json:"id"`
}

type DepositoEstoque struct {
	ID int64 `json:"id"`
}

// RespostaEstoqueBling representa a resposta da API ao atualizar estoque
type RespostaEstoqueBling struct {
	Data EstoqueData `json:"data"`
}

type EstoqueData struct {
	ID int64 `json:"id"`
}

// DepositoBling representa um depósito retornado pela API
type DepositoBling struct {
	ID                 int64  `json:"id"`
	Descricao          string `json:"descricao"`
	Situacao           int    `json:"situacao"` // 1 = Ativo
	Padrao             bool   `json:"padrao"`
	DesconsiderarSaldo bool   `json:"desconsiderarSaldo"`
}

// ListaDepositosBling representa a resposta da API de depósitos
type ListaDepositosBling struct {
	Data []DepositoBling `json:"data"`
}
