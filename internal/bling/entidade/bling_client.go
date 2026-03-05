package entidade

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	blingDTO "prosync/internal/bling/dto"
	"prosync/internal/comum/logger"
)

// BlingClient é o cliente para interagir com a API do Bling
type BlingClient struct {
	baseURL      string
	clientID     string
	clientSecret string
	accessToken  string
	httpClient   *http.Client
	logger       *logger.Logger
	lastRequest  time.Time // Para rate limiting
}

// NovoBlingClient cria uma nova instância do cliente Bling
func NovoBlingClient(clientID, clientSecret, accessToken string, logger *logger.Logger) *BlingClient {
	return &BlingClient{
		baseURL:      "https://api.bling.com.br/Api/v3",
		clientID:     clientID,
		clientSecret: clientSecret,
		accessToken:  accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger:      logger,
		lastRequest: time.Time{},
	}
}

// SetAccessToken atualiza o access token
func (c *BlingClient) SetAccessToken(token string) {
	c.accessToken = token
}

// rateLimitWait implementa rate limiting de 2 req/s (500ms entre requisições)
// Bling permite 3 req/s, mas usamos 2 req/s como margem de segurança
func (c *BlingClient) rateLimitWait() {
	if !c.lastRequest.IsZero() {
		elapsed := time.Since(c.lastRequest)
		minInterval := 500 * time.Millisecond // 2 req/s = 500ms entre requisições (margem de segurança)
		if elapsed < minInterval {
			time.Sleep(minInterval - elapsed)
		}
	}
	c.lastRequest = time.Now()
}

// doRequest executa uma requisição HTTP com rate limiting e logging
func (c *BlingClient) doRequest(method, endpoint string, body interface{}, response interface{}) error {
	c.rateLimitWait()

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("erro ao serializar payload: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, reqBody)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Trata rate limit (429)
	if resp.StatusCode == 429 {
		return &RateLimitError{
			Message:    "Rate limit excedido",
			RetryAfter: 300, // 5 minutos
		}
	}

	// Trata unauthorized (401) - token inválido
	if resp.StatusCode == 401 {
		return &UnauthorizedError{
			Message: "Token inválido ou expirado",
		}
	}

	// Trata outros erros HTTP
	if resp.StatusCode >= 400 {
		var errResp blingDTO.ErroAPIBling
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			return &APIError{
				StatusCode: resp.StatusCode,
				Message:    errResp.Error.Message,
				Type:       errResp.Error.Type,
			}
		}
		// Formata JSON de erro com acentuação correta
		return fmt.Errorf("%s", formatarJSONErro(respBody))
	}

	// Parse da resposta de sucesso
	if response != nil && resp.StatusCode < 300 {
		if err := json.Unmarshal(respBody, response); err != nil {
			return fmt.Errorf("erro ao decodificar resposta: %w", err)
		}
	}

	return nil
}

// RefreshToken renova o access token usando o refresh token
func (c *BlingClient) RefreshToken(refreshToken string) (*blingDTO.TokenRefreshResponse, error) {
	c.rateLimitWait()

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://www.bling.com.br/Api/v3/oauth/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição de refresh: %w", err)
	}

	// Basic Auth: base64(client_id:client_secret)
	auth := base64.StdEncoding.EncodeToString([]byte(c.clientID + ":" + c.clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar refresh token: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta de refresh: %w", err)
	}

	if resp.StatusCode == 429 {
		return nil, &RateLimitError{
			Message:    "Rate limit no refresh token",
			RetryAfter: 300,
		}
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("erro ao renovar token (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var tokenResp blingDTO.TokenRefreshResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta de refresh: %w", err)
	}

	// Atualiza o token do cliente
	c.accessToken = tokenResp.AccessToken

	return &tokenResp, nil
}

// BuscarProdutoPorCodigo busca um produto pelo SKU/código
func (c *BlingClient) BuscarProdutoPorCodigo(codigo string) (*blingDTO.ProdutoBlingData, error) {
	endpoint := fmt.Sprintf("/produtos?codigo=%s", url.QueryEscape(codigo))

	var resp blingDTO.ListaProdutosBling
	if err := c.doRequest("GET", endpoint, nil, &resp); err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, nil // Produto não encontrado
	}

	return &resp.Data[0], nil
}

// CriarProduto cria um novo produto no Bling
func (c *BlingClient) CriarProduto(produto *blingDTO.ProdutoBling) (*blingDTO.ProdutoBlingData, error) {
	var resp blingDTO.RespostaProdutoBling
	if err := c.doRequest("POST", "/produtos", produto, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// AtualizarProduto atualiza um produto existente no Bling
func (c *BlingClient) AtualizarProduto(id int64, produto *blingDTO.ProdutoBling) error {
	endpoint := fmt.Sprintf("/produtos/%d", id)
	return c.doRequest("PATCH", endpoint, produto, nil)
}

// BuscarDepositos lista todos os depósitos
func (c *BlingClient) BuscarDepositos() ([]blingDTO.DepositoBling, error) {
	var resp blingDTO.ListaDepositosBling
	if err := c.doRequest("GET", "/depositos", nil, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// AtualizarEstoque atualiza o estoque de um produto
func (c *BlingClient) AtualizarEstoque(estoque *blingDTO.EstoqueBling) error {
	var resp blingDTO.RespostaEstoqueBling
	return c.doRequest("POST", "/estoques", estoque, &resp)
}

// formatarJSONErro decodifica e formata JSON com acentuação correta
func formatarJSONErro(jsonBytes []byte) string {
	var data interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		// Se não conseguir decodificar, retorna string bruta
		return string(jsonBytes)
	}

	// Re-codifica com SetEscapeHTML(false) para preservar acentuação
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(data); err != nil {
		return string(jsonBytes)
	}

	// Remove quebra de linha extra que Encode adiciona
	result := buffer.String()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}

	return result
}

// Erros customizados
type RateLimitError struct {
	Message    string
	RetryAfter int // segundos
}

func (e *RateLimitError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

type APIError struct {
	StatusCode int
	Message    string
	Type       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Erro API Bling (HTTP %d): %s", e.StatusCode, e.Message)
}
