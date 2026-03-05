package dto

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexInt é um tipo que aceita int ou string no JSON e converte para int
type FlexInt int

// UnmarshalJSON implementa json.Unmarshaler para FlexInt
func (fi *FlexInt) UnmarshalJSON(data []byte) error {
	// Tenta fazer unmarshal como int primeiro
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*fi = FlexInt(i)
		return nil
	}

	// Se falhar, tenta como string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("FlexInt deve ser int ou string: %w", err)
	}

	// Converte string para int
	i, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("erro ao converter string para int: %w", err)
	}

	*fi = FlexInt(i)
	return nil
}

// Int retorna o valor como int
func (fi FlexInt) Int() int {
	return int(fi)
}

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
