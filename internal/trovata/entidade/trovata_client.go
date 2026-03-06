package entidade

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"prosync/internal/comum/logger"
	"prosync/internal/trovata/dto"
)

// TrovataClient é responsável por comunicação com a API Trovata
type TrovataClient struct {
	baseURL    string
	token      string
	empresaID  string
	httpClient *http.Client
	logger     *logger.Logger
}

// NovoTrovataClient cria nova instância do cliente Trovata
func NovoTrovataClient(logger *logger.Logger) *TrovataClient {
	return &TrovataClient{
		baseURL:   "https://geek-api.trovata.com.br/api_up",
		token:     "TCoKcewkvnosaQ7rJfkBUN81LqVEG",
		empresaID: "1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// CriarProduto envia produto para criar na Trovata
func (c *TrovataClient) CriarProduto(produto *dto.ProdutoTrovataRequest, sku string, idProdutoTiny string) error {
	inicio := time.Now()

	url := fmt.Sprintf("%s/produto/%s/%s", c.baseURL, c.token, c.empresaID)

	// Converte produto para array (API Trovata espera array)
	payload := []dto.ProdutoTrovataRequest{*produto}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar produto: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Desserializa para log completo
	var requisicaoCompleta interface{}
	json.Unmarshal(jsonData, &requisicaoCompleta)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		duracao := time.Since(inicio)
		c.logger.RegistrarChamada(logger.EntradaLog{
			Servico:       "trovata",
			Operacao:      "CriarProdutoTrovata",
			URL:           url,
			MetodoHTTP:    "POST",
			SKU:           sku,
			ProdutoTinyID: idProdutoTiny,
			Requisicao:    requisicaoCompleta,
			Erro:          err.Error(),
			Duracao:       duracao.String(),
			DuracaoMs:     float64(duracao.Milliseconds()),
		})
		return fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Parse da resposta
	var respostaCompleta interface{}
	json.Unmarshal(bodyBytes, &respostaCompleta)

	duracao := time.Since(inicio)
	c.logger.RegistrarChamada(logger.EntradaLog{
		Servico:       "trovata",
		Operacao:      "CriarProdutoTrovata",
		URL:           url,
		MetodoHTTP:    "POST",
		SKU:           sku,
		ProdutoTinyID: idProdutoTiny,
		Requisicao:    requisicaoCompleta,
		StatusCode:    resp.StatusCode,
		Resposta:      respostaCompleta,
		Duracao:       duracao.String(),
		DuracaoMs:     float64(duracao.Milliseconds()),
	})

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("erro HTTP %d ao criar produto na Trovata", resp.StatusCode)
	}

	return nil
}

// AtualizarEstoque atualiza o estoque de um produto na Trovata
func (c *TrovataClient) AtualizarEstoque(estoque *dto.EstoqueTrovataRequest, sku string, idProdutoTiny string) error {
	inicio := time.Now()

	url := fmt.Sprintf("%s/saldoEstoque/%s/%s", c.baseURL, c.token, c.empresaID)

	// Converte estoque para array (API Trovata espera array)
	payload := []dto.EstoqueTrovataRequest{*estoque}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar estoque: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Desserializa para log completo
	var requisicaoCompleta interface{}
	json.Unmarshal(jsonData, &requisicaoCompleta)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		duracao := time.Since(inicio)
		c.logger.RegistrarChamada(logger.EntradaLog{
			Servico:       "trovata",
			Operacao:      "AtualizarEstoqueTrovata",
			URL:           url,
			MetodoHTTP:    "POST",
			SKU:           sku,
			ProdutoTinyID: idProdutoTiny,
			Requisicao:    requisicaoCompleta,
			Erro:          err.Error(),
			Duracao:       duracao.String(),
			DuracaoMs:     float64(duracao.Milliseconds()),
		})
		return fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Parse da resposta
	var respostaCompleta interface{}
	json.Unmarshal(bodyBytes, &respostaCompleta)

	duracao := time.Since(inicio)
	c.logger.RegistrarChamada(logger.EntradaLog{
		Servico:       "trovata",
		Operacao:      "AtualizarEstoqueTrovata",
		URL:           url,
		MetodoHTTP:    "POST",
		SKU:           sku,
		ProdutoTinyID: idProdutoTiny,
		Requisicao:    requisicaoCompleta,
		StatusCode:    resp.StatusCode,
		Resposta:      respostaCompleta,
		Duracao:       duracao.String(),
		DuracaoMs:     float64(duracao.Milliseconds()),
	})

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("erro HTTP %d ao atualizar estoque na Trovata", resp.StatusCode)
	}

	return nil
}

// AtualizarStatusPedido atualiza o status de um pedido na Trovata
func (c *TrovataClient) AtualizarStatusPedido(ocorrencia *dto.OcorrenciaVendaRequest, sku string, idProdutoTiny string) error {
	inicio := time.Now()

	url := fmt.Sprintf("%s/ocorrenciaVenda/%s/%s", c.baseURL, c.token, c.empresaID)

	// Converte ocorrencia para array (API Trovata espera array)
	payload := []dto.OcorrenciaVendaRequest{*ocorrencia}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar ocorrência: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Desserializa para log completo
	var requisicaoCompleta interface{}
	json.Unmarshal(jsonData, &requisicaoCompleta)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		duracao := time.Since(inicio)
		c.logger.RegistrarChamada(logger.EntradaLog{
			Servico:       "trovata",
			Operacao:      "AtualizarStatusPedidoTrovata",
			URL:           url,
			MetodoHTTP:    "POST",
			SKU:           sku,
			ProdutoTinyID: idProdutoTiny,
			Requisicao:    requisicaoCompleta,
			Erro:          err.Error(),
			Duracao:       duracao.String(),
			DuracaoMs:     float64(duracao.Milliseconds()),
		})
		return fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Parse da resposta
	var respostaCompleta interface{}
	json.Unmarshal(bodyBytes, &respostaCompleta)

	duracao := time.Since(inicio)
	c.logger.RegistrarChamada(logger.EntradaLog{
		Servico:       "trovata",
		Operacao:      "AtualizarStatusPedidoTrovata",
		URL:           url,
		MetodoHTTP:    "POST",
		SKU:           sku,
		ProdutoTinyID: idProdutoTiny,
		Requisicao:    requisicaoCompleta,
		StatusCode:    resp.StatusCode,
		Resposta:      respostaCompleta,
		Duracao:       duracao.String(),
		DuracaoMs:     float64(duracao.Milliseconds()),
	})

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("erro HTTP %d ao atualizar status do pedido na Trovata", resp.StatusCode)
	}

	return nil
}
