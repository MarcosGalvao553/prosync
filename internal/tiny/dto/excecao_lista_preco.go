package dto

// ExcecaoListaPrecoRequest representa os dados de requisição para buscar exceções de lista de preços
type ExcecaoListaPrecoRequest struct {
	Token        string `json:"token"`
	IdListaPreco string `json:"idListaPreco"`
	Formato      string `json:"formato"`
	Pagina       int    `json:"pagina,omitempty"`
}

// NovaExcecaoListaPrecoRequest cria uma nova requisição com valores padrão
func NovaExcecaoListaPrecoRequest(token, idListaPreco string, pagina int) *ExcecaoListaPrecoRequest {
	return &ExcecaoListaPrecoRequest{
		Token:        token,
		IdListaPreco: idListaPreco,
		Formato:      "json",
		Pagina:       pagina,
	}
}

// ExcecaoListaPrecoResponse representa a resposta completa da API
type ExcecaoListaPrecoResponse struct {
	Retorno RetornoExcecaoListaPreco `json:"retorno"`
}

// RetornoExcecaoListaPreco representa o objeto de retorno dentro da resposta
type RetornoExcecaoListaPreco struct {
	StatusProcessamento string                             `json:"status_processamento"`
	Status              string                             `json:"status"`
	Pagina              FlexInt                            `json:"pagina"`
	NumeroPaginas       FlexInt                            `json:"numero_paginas"`
	Registros           []RegistroExcecaoListaPrecoWrapper `json:"registros"`
}

// RegistroExcecaoListaPrecoWrapper encapsula um registro de exceção
type RegistroExcecaoListaPrecoWrapper struct {
	Registro RegistroExcecaoListaPreco `json:"registro"`
}

// RegistroExcecaoListaPreco representa um registro individual de exceção de preço
type RegistroExcecaoListaPreco struct {
	ID           int     `json:"id"`
	IdListaPreco int     `json:"id_lista_preco"`
	IdProduto    int64   `json:"id_produto"`
	Preco        float64 `json:"preco"`
}

// ProdutoExcecaoListaPrecoTiny é o objeto simplificado que será usado no processamento
type ProdutoExcecaoListaPrecoTiny struct {
	ID        int     `json:"id"`
	IdProduto int64   `json:"id_produto"`
	Preco     float64 `json:"preco"`
}

// ParaProdutoExcecaoListaPrecoTiny converte um registro completo para o formato simplificado
func (r *RegistroExcecaoListaPreco) ParaProdutoExcecaoListaPrecoTiny() ProdutoExcecaoListaPrecoTiny {
	return ProdutoExcecaoListaPrecoTiny{
		ID:        r.ID,
		IdProduto: r.IdProduto,
		Preco:     r.Preco,
	}
}
