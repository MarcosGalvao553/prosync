package servidor

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"prosync/internal/comum/database"
	"prosync/internal/comum/logger"
)

// ProcessadorInterface define a interface para processar produtos
type ProcessadorInterface interface {
	ProcessarProdutoEspecifico(idProduto string) error
}

// ProcessadorBlingInterface define a interface para processar produtos no Bling
type ProcessadorBlingInterface interface {
	SincronizarProdutoParaUsuario(productID uint64, userID uint64) error
}

// ServidorWeb gerencia o servidor HTTP
type ServidorWeb struct {
	porta            string
	logger           *logger.Logger
	processador      ProcessadorInterface
	processadorBling ProcessadorBlingInterface
}

// NovoServidorWeb cria uma nova instância do servidor
func NovoServidorWeb(porta string, logger *logger.Logger) *ServidorWeb {
	return &ServidorWeb{
		porta:  porta,
		logger: logger,
	}
}

// SetProcessador define o processador que será usado para processar produtos
func (s *ServidorWeb) SetProcessador(processador ProcessadorInterface) {
	s.processador = processador
}

// SetProcessadorBling define o processador Bling
func (s *ServidorWeb) SetProcessadorBling(processadorBling ProcessadorBlingInterface) {
	s.processadorBling = processadorBling
}

// Iniciar inicia o servidor HTTP
func (s *ServidorWeb) Iniciar() error {
	// Serve arquivos estáticos da pasta web
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	// API endpoints
	http.HandleFunc("/api/logs", s.handlerLogs)
	http.HandleFunc("/api/logs/estatisticas", s.handlerEstatisticas)
	http.HandleFunc("/api/logs/tempo-por-operacao", s.handlerTempoPorOperacao)
	http.HandleFunc("/api/users", s.handlerUsers)
	http.HandleFunc("/api/health", s.handlerHealth)
	http.HandleFunc("/api/process-product", s.handlerProcessProduct)
	http.HandleFunc("/api/create-bling-product", s.handlerCreateBlingProduct)

	endereco := fmt.Sprintf(":%s", s.porta)
	log.Printf("🌐 Servidor web iniciado em http://localhost%s", endereco)

	return http.ListenAndServe(endereco, nil)
}

// handlerHealth verifica se o servidor está funcionando
func (s *ServidorWeb) handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "OK",
		"time":   time.Now(),
	})
}

// handlerEstatisticas retorna estatísticas agregadas
func (s *ServidorWeb) handlerEstatisticas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	filtros := s.extrairFiltros(r)

	stats, err := database.BuscarEstatisticas(filtros)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

// handlerTempoPorOperacao retorna tempo médio por operação
func (s *ServidorWeb) handlerTempoPorOperacao(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	filtros := s.extrairFiltros(r)

	resultados, err := database.BuscarTempoPorOperacao(filtros)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resultados)
}

// handlerLogs retorna logs do banco de dados com filtros aplicados
func (s *ServidorWeb) handlerLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	filtros := s.extrairFiltros(r)

	// Busca logs do banco
	logsDB, err := database.BuscarLogs(filtros, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Converte para formato JSON do dashboard
	logs := []map[string]interface{}{}
	for _, log := range logsDB {
		logs = append(logs, database.ConverterLogParaJSON(log))
	}

	json.NewEncoder(w).Encode(logs)
}

// handlerUsers retorna usuários para autocomplete
func (s *ServidorWeb) handlerUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	termo := r.URL.Query().Get("q")

	usuarios, err := database.BuscarUsuarios(termo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(usuarios)
}

// extrairFiltros extrai filtros dos parâmetros da query
func (s *ServidorWeb) extrairFiltros(r *http.Request) map[string]string {
	filtros := make(map[string]string)

	if servico := r.URL.Query().Get("servico"); servico != "" {
		filtros["servico"] = servico
	}

	if operacao := r.URL.Query().Get("operacao"); operacao != "" {
		filtros["operacao"] = operacao
	}

	if status := r.URL.Query().Get("status"); status != "" {
		filtros["status"] = status
	}

	if produtoID := r.URL.Query().Get("produto"); produtoID != "" {
		filtros["produto_tiny_id"] = produtoID
	}

	if sku := r.URL.Query().Get("sku"); sku != "" {
		filtros["sku"] = sku
	}

	if dataInicio := r.URL.Query().Get("data_inicio"); dataInicio != "" {
		filtros["data_inicio"] = dataInicio
	}

	if dataFim := r.URL.Query().Get("data_fim"); dataFim != "" {
		filtros["data_fim"] = dataFim
	}

	if userID := r.URL.Query().Get("user_id"); userID != "" {
		filtros["user_id"] = userID
	}

	return filtros
}

// handlerProcessProduct processa um produto específico (chamado pelo webhook)
func (s *ServidorWeb) handlerProcessProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Apenas POST é permitido
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verifica se o processador está configurado
	if s.processador == nil {
		http.Error(w, "Processador não configurado", http.StatusInternalServerError)
		return
	}

	// Lê o body da requisição (JSON do webhook Tiny)
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler body da requisição", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse do JSON do webhook
	var webhookData map[string]interface{}
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &webhookData); err != nil {
			s.logger.RegistrarInfo("webhook", fmt.Sprintf("Body não é JSON válido, tentando extrair da URL: %v", err))
		}
	}

	// Extrai o ID do produto
	var idProduto string

	// Tenta extrair do JSON em múltiplas estruturas possíveis
	if webhookData != nil {
		// 1. Tenta extrair de idProdutoTiny (campo direto)
		if id, ok := webhookData["idProdutoTiny"].(string); ok {
			idProduto = id
		} else if id, ok := webhookData["idProdutoTiny"].(float64); ok {
			idProduto = fmt.Sprintf("%.0f", id)
		}

		// 2. Se não encontrou, tenta extrair de data.dados.idProduto
		if idProduto == "" {
			if data, ok := webhookData["data"].(map[string]interface{}); ok {
				if dados, ok := data["dados"].(map[string]interface{}); ok {
					if id, ok := dados["idProduto"].(string); ok {
						idProduto = id
					} else if id, ok := dados["idProduto"].(float64); ok {
						idProduto = fmt.Sprintf("%.0f", id)
					}
				}
			}
		}

		// 3. Se não encontrou, tenta extrair de dados.idProduto (estrutura original)
		if idProduto == "" {
			if dados, ok := webhookData["dados"].(map[string]interface{}); ok {
				if id, ok := dados["idProduto"].(string); ok {
					idProduto = id
				} else if id, ok := dados["idProduto"].(float64); ok {
					idProduto = fmt.Sprintf("%.0f", id)
				}
			}
		}
	}

	// Se não encontrou no JSON, tenta extrair da URL
	if idProduto == "" {
		path := r.URL.Path
		parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
		if len(parts) >= 4 {
			idProduto = parts[len(parts)-1]
		}
	}

	// Valida se conseguiu extrair o ID
	if idProduto == "" {
		http.Error(w, "ID do produto não informado (nem no JSON nem na URL)", http.StatusBadRequest)
		return
	}

	// Registra a notificação do Tiny no monitoramento
	s.logger.RegistrarChamada(logger.EntradaLog{
		Servico:       "tiny",
		Operacao:      "NotificacaoEstoqueTiny",
		ProdutoTinyID: idProduto,
		Requisicao:    webhookData,
		MetodoHTTP:    r.Method,
		URL:           r.URL.String(),
	})

	s.logger.RegistrarInfo("webhook", fmt.Sprintf("Webhook Tiny recebido - Produto %s", idProduto))

	// Processa o produto em BACKGROUND (não bloqueia a resposta)
	// Responde imediatamente e processa de forma assíncrona
	go func() {
		if err := s.processador.ProcessarProdutoEspecifico(idProduto); err != nil {
			s.logger.RegistrarErro("webhook", fmt.Sprintf("Erro ao processar produto %s", idProduto), err)
		} else {
			s.logger.RegistrarInfo("webhook", fmt.Sprintf("Produto %s processado com sucesso em background", idProduto))
		}
	}()

	// Retorna sucesso IMEDIATAMENTE
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Produto %s adicionado à fila de processamento", idProduto),
		"id":      idProduto,
	})
}

// handlerCreateBlingProduct cria um produto no Bling para um usuário específico
func (s *ServidorWeb) handlerCreateBlingProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Apenas POST é permitido
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verifica se o processador Bling está configurado
	if s.processadorBling == nil {
		http.Error(w, "Processador Bling não configurado", http.StatusInternalServerError)
		return
	}

	// Lê o body da requisição
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler body da requisição", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse do JSON
	var requestData struct {
		ProductID uint64 `json:"product_id"`
		UserID    uint64 `json:"user_id"`
	}

	if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Valida os parâmetros
	if requestData.ProductID == 0 {
		http.Error(w, "product_id é obrigatório", http.StatusBadRequest)
		return
	}

	if requestData.UserID == 0 {
		http.Error(w, "user_id é obrigatório", http.StatusBadRequest)
		return
	}

	// Registra a requisição no log
	s.logger.RegistrarChamada(logger.EntradaLog{
		Servico:    "bling",
		Operacao:   "RequisicaoCriarProdutoBling",
		UserID:     requestData.UserID,
		Requisicao: requestData,
		MetodoHTTP: r.Method,
		URL:        r.URL.String(),
	})

	s.logger.RegistrarInfo("api", fmt.Sprintf("Requisição para criar produto %d no Bling do usuário %d", requestData.ProductID, requestData.UserID))

	// Processa em BACKGROUND (não bloqueia a resposta)
	// Responde imediatamente e processa de forma assíncrona
	go func() {
		if err := s.processadorBling.SincronizarProdutoParaUsuario(requestData.ProductID, requestData.UserID); err != nil {
			s.logger.RegistrarErro("bling",
				fmt.Sprintf("Erro ao criar produto %d no Bling para usuário %d", requestData.ProductID, requestData.UserID),
				err,
			)
		} else {
			s.logger.RegistrarInfo("bling",
				fmt.Sprintf("Produto %d criado/atualizado com sucesso no Bling para usuário %d", requestData.ProductID, requestData.UserID),
			)
		}
	}()

	// Retorna sucesso IMEDIATAMENTE
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"message":    fmt.Sprintf("Produto %d adicionado à fila de processamento do Bling para usuário %d", requestData.ProductID, requestData.UserID),
		"product_id": requestData.ProductID,
		"user_id":    requestData.UserID,
	})
}
