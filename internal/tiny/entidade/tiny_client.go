package entidade

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"prosync/internal/comum/config"
	"prosync/internal/comum/logger"
	"prosync/internal/tiny/dto"
)

// TinyClient é a entidade responsável por fazer chamadas à API do Tiny
type TinyClient struct {
	config           *config.Config
	logger           *logger.Logger
	httpClient       *http.Client
	ultimaRequisicao time.Time
}

// NovoTinyClient cria uma nova instância do cliente Tiny
func NovoTinyClient(cfg *config.Config, log *logger.Logger) *TinyClient {
	return &TinyClient{
		config: cfg,
		logger: log,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ObterConfig retorna a configuração do cliente
func (t *TinyClient) ObterConfig() *config.Config {
	return t.config
}

// BuscarExcecoesListaPreco busca as exceções de lista de preços de uma página específica
func (t *TinyClient) BuscarExcecoesListaPreco(pagina int) (*dto.ExcecaoListaPrecoResponse, error) {
	inicio := time.Now()

	// Monta a URL
	urlCompleta := fmt.Sprintf("%s/listas.precos.excecoes.php", t.config.TinyBaseURL)

	// Prepara os dados do formulário
	formData := url.Values{}
	formData.Set("token", t.config.TinyBearerToken)
	formData.Set("idListaPreco", t.config.TinyIdListaPreco)
	formData.Set("formato", "json")
	formData.Set("idProduto", "1025799210")
	if pagina > 0 {
		formData.Set("pagina", fmt.Sprintf("%d", pagina))
	}

	// Cria a requisição
	req, err := http.NewRequest("POST", urlCompleta, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Define os headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.config.TinyBearerToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Prepara dados para log (sem expor o token completo)
	dadosRequisicao := map[string]interface{}{
		"idListaPreco": t.config.TinyIdListaPreco,
		"formato":      "json",
		"pagina":       pagina,
		"token":        "***OCULTO***",
	}

	// Faz a requisição
	resp, err := t.httpClient.Do(req)
	if err != nil {
		// Registra erro no log
		duracao := time.Since(inicio)
		t.logger.RegistrarChamada(logger.EntradaLog{
			Servico:    "tiny",
			Operacao:   "BuscarExcecoesListaPreco",
			URL:        urlCompleta,
			MetodoHTTP: "POST",
			Requisicao: dadosRequisicao,
			Erro:       err.Error(),
			Duracao:    duracao.String(),
			DuracaoMs:  float64(duracao.Milliseconds()),
		})
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// Registra erro no log
		duracao := time.Since(inicio)
		t.logger.RegistrarChamada(logger.EntradaLog{
			Servico:    "tiny",
			Operacao:   "BuscarExcecoesListaPreco",
			URL:        urlCompleta,
			MetodoHTTP: "POST",
			Requisicao: dadosRequisicao,
			StatusCode: resp.StatusCode,
			Erro:       fmt.Sprintf("erro ao ler resposta: %v", err),
			Duracao:    duracao.String(),
			DuracaoMs:  float64(duracao.Milliseconds()),
		})
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Parse da resposta completa para o log
	var respostaCompleta map[string]interface{}
	json.Unmarshal(bodyBytes, &respostaCompleta)

	// Parse do JSON para o DTO
	var resposta dto.ExcecaoListaPrecoResponse
	if err := json.Unmarshal(bodyBytes, &resposta); err != nil {
		// Registra erro no log
		duracao := time.Since(inicio)
		t.logger.RegistrarChamada(logger.EntradaLog{
			Servico:    "tiny",
			Operacao:   "BuscarExcecoesListaPreco",
			URL:        urlCompleta,
			MetodoHTTP: "POST",
			Requisicao: dadosRequisicao,
			StatusCode: resp.StatusCode,
			Resposta:   respostaCompleta,
			Erro:       fmt.Sprintf("erro ao fazer parse do JSON: %v", err),
			Duracao:    duracao.String(),
			DuracaoMs:  float64(duracao.Milliseconds()),
		})
		return nil, fmt.Errorf("erro ao fazer parse do JSON: %w", err)
	}

	// Registra chamada bem-sucedida no log com resposta completa
	duracao := time.Since(inicio)
	t.logger.RegistrarChamada(logger.EntradaLog{
		Servico:    "tiny",
		Operacao:   "BuscarExcecoesListaPreco",
		URL:        urlCompleta,
		MetodoHTTP: "POST",
		Requisicao: dadosRequisicao,
		Resposta:   respostaCompleta,
		StatusCode: resp.StatusCode,
		Duracao:    duracao.String(),
		DuracaoMs:  float64(duracao.Milliseconds()),
	})

	return &resposta, nil
}

// BuscarTodasExcecoesListaPreco busca todas as exceções de lista de preços, iterando por todas as páginas
func (t *TinyClient) BuscarTodasExcecoesListaPreco() ([]dto.ProdutoExcecaoListaPrecoTiny, error) {
	t.logger.RegistrarInfo("tiny", "Iniciando busca de todas as exceções de lista de preços")

	var todosProdutos []dto.ProdutoExcecaoListaPrecoTiny
	paginaAtual := 1

	for {
		// Busca a página atual
		resposta, err := t.BuscarExcecoesListaPreco(paginaAtual)
		if err != nil {
			t.logger.RegistrarErro("tiny", fmt.Sprintf("Erro ao buscar página %d", paginaAtual), err)
			return nil, err
		}

		// Verifica se a API retornou erro (status diferente de OK)
		if resposta.Retorno.Status != "OK" {
			t.logger.RegistrarInfo("tiny", fmt.Sprintf("API retornou status '%s' na página %d - finalizando busca",
				resposta.Retorno.Status, paginaAtual))
			break
		}

		// Adiciona os produtos da página atual (limitado a 5 para teste)
		registrosProcessados := 0
		for _, wrapper := range resposta.Retorno.Registros {
			// if registrosProcessados >= 5 {
			// 	break
			// }
			produto := wrapper.Registro.ParaProdutoExcecaoListaPrecoTiny()
			todosProdutos = append(todosProdutos, produto)
			registrosProcessados++
		}

		t.logger.RegistrarInfo("tiny", fmt.Sprintf("Página %d/%d processada - %d registros",
			paginaAtual, resposta.Retorno.NumeroPaginas.Int(), registrosProcessados))

		// Verifica se há mais páginas
		if paginaAtual >= resposta.Retorno.NumeroPaginas.Int() {
			break
		}

		paginaAtual++

		// Aguarda um pouco entre requisições para não sobrecarregar a API
		time.Sleep(500 * time.Millisecond)
	}

	t.logger.RegistrarInfo("tiny", fmt.Sprintf("Busca concluída - Total de %d produtos", len(todosProdutos)))

	return todosProdutos, nil
}

// aguardarRateLimit garante que respeitamos o limite de 1 requisição por segundo
func (t *TinyClient) aguardarRateLimit() {
	agora := time.Now()
	tempoDecorrido := agora.Sub(t.ultimaRequisicao)

	// Se passou menos de 1 segundo desde a última requisição, aguarda
	if tempoDecorrido < time.Second {
		tempoEspera := time.Second - tempoDecorrido
		time.Sleep(tempoEspera)
	}

	t.ultimaRequisicao = time.Now()
}

// tratarRateLimit verifica se houve erro de rate limit e aguarda se necessário
func (t *TinyClient) tratarRateLimit(codigoErro string) bool {
	// Código 6 = API Bloqueada por excesso de requisições
	if codigoErro == "6" {
		t.logger.RegistrarInfo("tiny", "Rate limit atingido - aguardando 1 minuto")
		time.Sleep(1 * time.Minute)
		return true
	}
	return false
}

// BuscarDadosProduto busca os dados completos de um produto pelo ID
func (t *TinyClient) BuscarDadosProduto(idProduto string) (*dto.ProdutoTiny, error) {
	inicio := time.Now()

	// Aguarda rate limit
	t.aguardarRateLimit()

	// Monta a URL
	urlCompleta := fmt.Sprintf("%s/produto.obter.php", t.config.TinyBaseURL)

	// Prepara os dados do formulário
	formData := url.Values{}
	formData.Set("token", t.config.TinyBearerToken)
	formData.Set("id", idProduto)
	formData.Set("formato", "json")

	// Cria a requisição
	req, err := http.NewRequest("POST", urlCompleta, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Prepara dados para log
	dadosRequisicao := map[string]interface{}{
		"id":      idProduto,
		"formato": "json",
		"token":   "***OCULTO***",
	}

	// Faz a requisição
	resp, err := t.httpClient.Do(req)
	if err != nil {
		duracao := time.Since(inicio)
		t.logger.RegistrarChamada(logger.EntradaLog{
			Servico:       "tiny",
			Operacao:      "BuscarDadosProduto",
			URL:           urlCompleta,
			MetodoHTTP:    "POST",
			Requisicao:    dadosRequisicao,
			Erro:          err.Error(),
			Duracao:       duracao.String(),
			DuracaoMs:     float64(duracao.Milliseconds()),
			ProdutoTinyID: idProduto,
		})
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Verifica se a resposta é HTML (erro da API)
	if len(bodyBytes) > 0 && bodyBytes[0] == '<' {
		duracao := time.Since(inicio)
		t.logger.RegistrarChamada(logger.EntradaLog{
			Servico:       "tiny",
			Operacao:      "BuscarDadosProduto",
			URL:           urlCompleta,
			MetodoHTTP:    "POST",
			Requisicao:    dadosRequisicao,
			StatusCode:    resp.StatusCode,
			Erro:          fmt.Sprintf("API retornou HTML (status %d)", resp.StatusCode),
			Duracao:       duracao.String(),
			DuracaoMs:     float64(duracao.Milliseconds()),
			ProdutoTinyID: idProduto,
		})
		return nil, fmt.Errorf("API retornou HTML em vez de JSON (status %d) - possível erro de rate limit ou bloqueio", resp.StatusCode)
	}

	// Parse da resposta completa para o log
	var respostaCompleta map[string]interface{}
	json.Unmarshal(bodyBytes, &respostaCompleta)

	// Parse do JSON para o DTO
	var resposta dto.ProdutoResponse
	if err := json.Unmarshal(bodyBytes, &resposta); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do JSON: %w", err)
	}

	// Verifica rate limit
	if t.tratarRateLimit(resposta.Retorno.CodigoErro) {
		// Tenta novamente após aguardar
		return t.BuscarDadosProduto(idProduto)
	}

	// Verifica se houve erro
	if resposta.Retorno.Status != "OK" {
		if len(resposta.Retorno.Erros) > 0 {
			return nil, fmt.Errorf("erro da API: %s", resposta.Retorno.Erros[0].Erro)
		}
		return nil, fmt.Errorf("erro desconhecido da API")
	}

	produtoTiny := resposta.Retorno.Produto.ParaProdutoTiny()

	// Registra log com resposta completa (após parse para pegar SKU)
	duracao := time.Since(inicio)
	t.logger.RegistrarChamada(logger.EntradaLog{
		Servico:       "tiny",
		Operacao:      "BuscarDadosProduto",
		URL:           urlCompleta,
		MetodoHTTP:    "POST",
		Requisicao:    dadosRequisicao,
		Resposta:      respostaCompleta,
		StatusCode:    resp.StatusCode,
		Duracao:       duracao.String(),
		DuracaoMs:     float64(duracao.Milliseconds()),
		ProdutoTinyID: idProduto,
		SKU:           produtoTiny.Codigo,
	})

	return &produtoTiny, nil
}

// BuscarEstoqueProduto busca o estoque de um produto pelo ID
func (t *TinyClient) BuscarEstoqueProduto(idProduto string) (*dto.EstoqueTiny, error) {
	inicio := time.Now()

	// Aguarda rate limit
	t.aguardarRateLimit()

	// Monta a URL
	urlCompleta := fmt.Sprintf("%s/produto.obter.estoque.php", t.config.TinyBaseURL)

	// Prepara os dados do formulário
	formData := url.Values{}
	formData.Set("token", t.config.TinyBearerToken)
	formData.Set("id", idProduto)
	formData.Set("idListaPreco", t.config.TinyIdListaPreco)
	formData.Set("formato", "json")

	// Cria a requisição
	req, err := http.NewRequest("POST", urlCompleta, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.config.TinyBearerToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Prepara dados para log
	dadosRequisicao := map[string]interface{}{
		"id":           idProduto,
		"idListaPreco": t.config.TinyIdListaPreco,
		"token":        "***OCULTO***",
	}

	// Faz a requisição
	resp, err := t.httpClient.Do(req)
	if err != nil {
		duracao := time.Since(inicio)
		t.logger.RegistrarChamada(logger.EntradaLog{
			Servico:       "tiny",
			Operacao:      "BuscarEstoqueProduto",
			URL:           urlCompleta,
			MetodoHTTP:    "POST",
			Requisicao:    dadosRequisicao,
			Erro:          err.Error(),
			Duracao:       duracao.String(),
			DuracaoMs:     float64(duracao.Milliseconds()),
			ProdutoTinyID: idProduto,
		})
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Verifica se a resposta é HTML (erro da API)
	if len(bodyBytes) > 0 && bodyBytes[0] == '<' {
		duracao := time.Since(inicio)
		t.logger.RegistrarChamada(logger.EntradaLog{
			Servico:       "tiny",
			Operacao:      "BuscarEstoqueProduto",
			URL:           urlCompleta,
			MetodoHTTP:    "POST",
			Requisicao:    dadosRequisicao,
			StatusCode:    resp.StatusCode,
			Erro:          fmt.Sprintf("API retornou HTML (status %d)", resp.StatusCode),
			Duracao:       duracao.String(),
			DuracaoMs:     float64(duracao.Milliseconds()),
			ProdutoTinyID: idProduto,
		})
		return nil, fmt.Errorf("API retornou HTML em vez de JSON (status %d) - possível erro de rate limit ou bloqueio", resp.StatusCode)
	}

	// Parse da resposta completa para o log
	var respostaCompleta map[string]interface{}
	json.Unmarshal(bodyBytes, &respostaCompleta)

	// Parse do JSON para o DTO
	var resposta dto.EstoqueResponse
	if err := json.Unmarshal(bodyBytes, &resposta); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do JSON: %w", err)
	}

	// Verifica rate limit
	if t.tratarRateLimit(resposta.Retorno.CodigoErro) {
		// Tenta novamente após aguardar
		return t.BuscarEstoqueProduto(idProduto)
	}

	// Verifica se houve erro
	if resposta.Retorno.Status != "OK" {
		if len(resposta.Retorno.Erros) > 0 {
			return nil, fmt.Errorf("erro da API: %s", resposta.Retorno.Erros[0].Erro)
		}
		return nil, fmt.Errorf("erro desconhecido da API")
	}

	estoqueTiny := resposta.Retorno.Produto.ParaEstoqueTiny()

	// Registra log com resposta completa (após parse para pegar SKU)
	duracao := time.Since(inicio)
	t.logger.RegistrarChamada(logger.EntradaLog{
		Servico:       "tiny",
		Operacao:      "BuscarEstoqueProduto",
		URL:           urlCompleta,
		MetodoHTTP:    "POST",
		Requisicao:    dadosRequisicao,
		Resposta:      respostaCompleta,
		StatusCode:    resp.StatusCode,
		Duracao:       duracao.String(),
		DuracaoMs:     float64(duracao.Milliseconds()),
		ProdutoTinyID: idProduto,
		SKU:           estoqueTiny.Codigo,
	})

	return &estoqueTiny, nil
}

// BuscarPrecoProdutoListaPreco busca o preço de um produto específico na lista de preços
func (t *TinyClient) BuscarPrecoProdutoListaPreco(idListaPreco int, idProduto string) (*dto.ProdutoExcecaoListaPrecoTiny, error) {
	inicio := time.Now()

	// Monta a URL
	urlCompleta := fmt.Sprintf("%s/lista.preco.obter.php", t.config.TinyBaseURL)

	// Prepara os dados do formulário
	formData := url.Values{}
	formData.Set("token", t.config.TinyBearerToken)
	formData.Set("id", fmt.Sprintf("%d", idListaPreco))
	formData.Set("idProduto", idProduto)
	formData.Set("formato", "json")

	// Cria a requisição
	req, err := http.NewRequest("POST", urlCompleta, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Define os headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.config.TinyBearerToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Prepara dados para log (sem expor o token completo)
	dadosRequisicao := map[string]interface{}{
		"id":        idListaPreco,
		"idProduto": idProduto,
		"formato":   "json",
		"token":     "***OCULTO***",
	}

	// Faz a requisição
	resp, err := t.httpClient.Do(req)
	if err != nil {
		// Registra erro no log
		duracao := time.Since(inicio)
		t.logger.RegistrarChamada(logger.EntradaLog{
			Servico:       "tiny",
			Operacao:      "BuscarPrecoProdutoListaPreco",
			URL:           urlCompleta,
			MetodoHTTP:    "POST",
			Requisicao:    dadosRequisicao,
			Erro:          err.Error(),
			Duracao:       duracao.String(),
			DuracaoMs:     float64(duracao.Milliseconds()),
			ProdutoTinyID: idProduto,
		})
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// Registra erro no log
		duracao := time.Since(inicio)
		t.logger.RegistrarChamada(logger.EntradaLog{
			Servico:       "tiny",
			Operacao:      "BuscarPrecoProdutoListaPreco",
			URL:           urlCompleta,
			MetodoHTTP:    "POST",
			Requisicao:    dadosRequisicao,
			StatusCode:    resp.StatusCode,
			Erro:          err.Error(),
			Duracao:       duracao.String(),
			DuracaoMs:     float64(duracao.Milliseconds()),
			ProdutoTinyID: idProduto,
		})
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Parse da resposta completa para o log
	var respostaCompleta map[string]interface{}
	json.Unmarshal(bodyBytes, &respostaCompleta)

	// Parse do JSON para o DTO
	var resposta dto.ExcecaoListaPrecoResponse
	if err := json.Unmarshal(bodyBytes, &resposta); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do JSON: %w", err)
	}

	// Registra log com resposta completa
	duracao := time.Since(inicio)
	t.logger.RegistrarChamada(logger.EntradaLog{
		Servico:       "tiny",
		Operacao:      "BuscarPrecoProdutoListaPreco",
		URL:           urlCompleta,
		MetodoHTTP:    "POST",
		Requisicao:    dadosRequisicao,
		Resposta:      respostaCompleta,
		StatusCode:    resp.StatusCode,
		Duracao:       duracao.String(),
		DuracaoMs:     float64(duracao.Milliseconds()),
		ProdutoTinyID: idProduto,
	})

	// Verifica se houve erro
	if resposta.Retorno.Status != "OK" {
		return nil, fmt.Errorf("produto não encontrado na lista de preços")
	}

	// Verifica se há registros
	if len(resposta.Retorno.Registros) == 0 {
		return nil, fmt.Errorf("produto sem preço definido na lista de preços")
	}

	// Retorna o primeiro registro encontrado
	precoTiny := resposta.Retorno.Registros[0].Registro.ParaProdutoExcecaoListaPrecoTiny()
	return &precoTiny, nil
}
