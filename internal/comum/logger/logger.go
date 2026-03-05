package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"prosync/internal/comum/database"
)

type TipoLog string

const (
	TipoLogTiny  TipoLog = "tiny"
	TipoLogBling TipoLog = "bling"
)

type EntradaLog struct {
	Timestamp       time.Time              `json:"timestamp"`
	Servico         string                 `json:"servico"`
	Operacao        string                 `json:"operacao"`
	Requisicao      map[string]interface{} `json:"requisicao"`
	Resposta        map[string]interface{} `json:"resposta"`
	StatusCode      int                    `json:"status_code,omitempty"`
	Erro            string                 `json:"erro,omitempty"`
	Duracao         string                 `json:"duracao,omitempty"`
	DuracaoMs       float64                `json:"duracao_ms,omitempty"`
	URL             string                 `json:"url,omitempty"`
	MetodoHTTP      string                 `json:"metodo_http,omitempty"`
	ProdutoTinyID   string                 `json:"produto_tiny_id,omitempty"`
	SKU             string                 `json:"sku,omitempty"`
	IDProdutoBling  string                 `json:"idprodutobling,omitempty"`
	UserID          uint64                 `json:"user_id,omitempty"`
	RequestHeaders  map[string]interface{} `json:"request_headers,omitempty"`
	ResponseHeaders map[string]interface{} `json:"response_headers,omitempty"`
}

type Logger struct {
	pastLogs string
}

// NovoLogger cria uma nova instância do logger
func NovoLogger() (*Logger, error) {
	pastLogs := "logs"

	// Cria a pasta de logs se não existir
	if err := os.MkdirAll(pastLogs, 0755); err != nil {
		return nil, fmt.Errorf("erro ao criar pasta de logs: %w", err)
	}

	return &Logger{
		pastLogs: pastLogs,
	}, nil
}

// marshalJSON serializa para JSON preservando acentuação (sem escapar Unicode)
func marshalJSON(v interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(v); err != nil {
		return nil, err
	}

	// Remove quebra de linha que Encode adiciona
	result := buffer.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}

	return result, nil
}

// marshalIndentJSON serializa para JSON formatado preservando acentuação
func marshalIndentJSON(v interface{}, prefix, indent string) ([]byte, error) {
	data, err := marshalJSON(v)
	if err != nil {
		return nil, err
	}

	// Re-formata com indentação
	var formatted interface{}
	if err := json.Unmarshal(data, &formatted); err != nil {
		return data, nil
	}

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent(prefix, indent)

	if err := encoder.Encode(formatted); err != nil {
		return data, nil
	}

	// Remove quebra de linha que Encode adiciona
	result := buffer.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}

	return result, nil
}

// RegistrarChamada registra uma chamada HTTP no banco de dados (fallback para arquivo em caso de erro)
func (l *Logger) RegistrarChamada(entrada EntradaLog) error {
	entrada.Timestamp = time.Now()

	// Tenta salvar no banco de dados primeiro
	if err := l.salvarNoBanco(entrada); err != nil {
		// Se falhar, salva em arquivo como fallback
		fmt.Printf("⚠️  Erro ao salvar log no banco, usando arquivo: %v\n", err)
		return l.salvarEmArquivo(entrada)
	}

	return nil
}

// salvarNoBanco salva o log no banco de dados
func (l *Logger) salvarNoBanco(entrada EntradaLog) error {
	// Determina o status
	status := "OK"
	if entrada.Erro != "" {
		status = "Erro"
	} else if resposta, ok := entrada.Resposta["status"].(string); ok {
		if resposta == "Erro" {
			status = "Erro"
		}
	}

	// Extrai código de erro se existir
	errorCode := ""
	if resposta, ok := entrada.Resposta["codigo"].(string); ok {
		errorCode = resposta
	} else if resposta, ok := entrada.Resposta["codigo"].(float64); ok {
		errorCode = fmt.Sprintf("%.0f", resposta)
	}

	// Extrai produto_tiny_id da requisição ou resposta
	produtoTinyID := entrada.ProdutoTinyID
	if produtoTinyID == "" {
		if id, ok := entrada.Requisicao["id"].(string); ok {
			produtoTinyID = id
		} else if id, ok := entrada.Resposta["id_produto"].(string); ok {
			produtoTinyID = id
		}
	}

	// Serializa requisição e resposta (preservando acentuação)
	reqBody, _ := marshalJSON(entrada.Requisicao)
	respBody, _ := marshalJSON(entrada.Resposta)
	reqHeaders, _ := marshalJSON(entrada.RequestHeaders)
	respHeaders, _ := marshalJSON(entrada.ResponseHeaders)

	// Cria metadata com informações extras
	metadata := map[string]interface{}{}
	if entrada.Duracao != "" {
		metadata["duracao_formatada"] = entrada.Duracao
	}
	metadataJSON, _ := marshalJSON(metadata)

	log := &database.LogAPI{
		CreatedAt:          entrada.Timestamp,
		Servico:            entrada.Servico,
		Operacao:           entrada.Operacao,
		Status:             status,
		RequestMethod:      entrada.MetodoHTTP,
		RequestURL:         entrada.URL,
		RequestHeaders:     string(reqHeaders),
		RequestBody:        string(reqBody),
		RequestSizeBytes:   len(reqBody),
		ResponseStatusCode: entrada.StatusCode,
		ResponseHeaders:    string(respHeaders),
		ResponseBody:       string(respBody),
		ResponseSizeBytes:  len(respBody),
		DurationMs:         entrada.DuracaoMs,
		ProdutoTinyID:      produtoTinyID,
		SKU:                entrada.SKU,
		IDProdutoBling:     entrada.IDProdutoBling,
		UserID:             entrada.UserID,
		ErrorCode:          errorCode,
		ErrorMessage:       entrada.Erro,
		Metadata:           string(metadataJSON),
	}

	return database.SalvarLog(log)
}

// salvarEmArquivo salva o log em arquivos JSON e texto (fallback)
func (l *Logger) salvarEmArquivo(entrada EntradaLog) error {
	// Gera os nomes dos arquivos baseados no serviço e data
	data := entrada.Timestamp.Format("2006-01-02")
	nomeArquivoJSON := filepath.Join(l.pastLogs, fmt.Sprintf("%s_%s.json", entrada.Servico, data))
	nomeArquivoTexto := filepath.Join(l.pastLogs, fmt.Sprintf("%s_%s.log", entrada.Servico, data))

	// Registra em JSON
	if err := l.registrarJSON(nomeArquivoJSON, entrada); err != nil {
		return err
	}

	// Registra em texto legível
	if err := l.registrarTexto(nomeArquivoTexto, entrada); err != nil {
		return err
	}

	return nil
}

// registrarJSON salva o log em formato JSON
func (l *Logger) registrarJSON(nomeArquivo string, entrada EntradaLog) error {
	arquivo, err := os.OpenFile(nomeArquivo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo JSON: %w", err)
	}
	defer arquivo.Close()

	jsonData, err := marshalJSON(entrada)
	if err != nil {
		return fmt.Errorf("erro ao serializar JSON: %w", err)
	}

	if _, err := arquivo.WriteString(string(jsonData) + "\n"); err != nil {
		return fmt.Errorf("erro ao escrever no arquivo JSON: %w", err)
	}

	return nil
}

// registrarTexto salva o log em formato texto legível
func (l *Logger) registrarTexto(nomeArquivo string, entrada EntradaLog) error {
	arquivo, err := os.OpenFile(nomeArquivo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de texto: %w", err)
	}
	defer arquivo.Close()

	// Formata o log de forma legível
	linhasSeparador := "========================================\n"
	texto := fmt.Sprintf(
		"%s[%s] %s - %s\nURL: %s %s\nStatus: %d\n",
		linhasSeparador,
		entrada.Timestamp.Format("2006-01-02 15:04:05"),
		entrada.Servico,
		entrada.Operacao,
		entrada.MetodoHTTP,
		entrada.URL,
		entrada.StatusCode,
	)

	if entrada.Duracao != "" {
		texto += fmt.Sprintf("Duração: %s\n", entrada.Duracao)
	}

	// Requisição
	if len(entrada.Requisicao) > 0 {
		texto += "\n--- REQUISIÇÃO ---\n"
		reqJSON, _ := marshalIndentJSON(entrada.Requisicao, "", "  ")
		texto += string(reqJSON) + "\n"
	}

	// Resposta
	if len(entrada.Resposta) > 0 {
		texto += "\n--- RESPOSTA ---\n"
		respJSON, _ := marshalIndentJSON(entrada.Resposta, "", "  ")
		texto += string(respJSON) + "\n"
	}

	// Erro
	if entrada.Erro != "" {
		texto += fmt.Sprintf("\n!!! ERRO: %s\n", entrada.Erro)
	}

	texto += linhasSeparador + "\n"

	if _, err := arquivo.WriteString(texto); err != nil {
		return fmt.Errorf("erro ao escrever no arquivo de texto: %w", err)
	}

	return nil
}

// RegistrarInfo registra uma mensagem informativa
func (l *Logger) RegistrarInfo(servico, mensagem string) error {
	entrada := EntradaLog{
		Timestamp: time.Now(),
		Servico:   servico,
		Operacao:  "INFO",
		Resposta: map[string]interface{}{
			"mensagem": mensagem,
		},
	}

	data := entrada.Timestamp.Format("2006-01-02")
	nomeArquivoTexto := filepath.Join(l.pastLogs, fmt.Sprintf("%s_%s.log", servico, data))

	arquivo, err := os.OpenFile(nomeArquivoTexto, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer arquivo.Close()

	texto := fmt.Sprintf("[%s] INFO: %s\n", entrada.Timestamp.Format("2006-01-02 15:04:05"), mensagem)
	_, err = arquivo.WriteString(texto)
	return err
}

// RegistrarErro registra uma mensagem de erro
func (l *Logger) RegistrarErro(servico, mensagem string, erro error) error {
	entrada := EntradaLog{
		Timestamp: time.Now(),
		Servico:   servico,
		Operacao:  "ERRO",
		Erro:      erro.Error(),
		Resposta: map[string]interface{}{
			"mensagem": mensagem,
		},
	}

	data := entrada.Timestamp.Format("2006-01-02")
	nomeArquivoTexto := filepath.Join(l.pastLogs, fmt.Sprintf("%s_%s.log", servico, data))

	arquivo, err := os.OpenFile(nomeArquivoTexto, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer arquivo.Close()

	texto := fmt.Sprintf("[%s] ERRO: %s - %v\n", entrada.Timestamp.Format("2006-01-02 15:04:05"), mensagem, erro)
	_, err = arquivo.WriteString(texto)
	return err
}
