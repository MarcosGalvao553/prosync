package servidor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"prosync/internal/comum/database"
	"prosync/internal/comum/logger"
)

// ServidorWeb gerencia o servidor HTTP
type ServidorWeb struct {
	porta  string
	logger *logger.Logger
}

// NovoServidorWeb cria uma nova instância do servidor
func NovoServidorWeb(porta string, logger *logger.Logger) *ServidorWeb {
	return &ServidorWeb{
		porta:  porta,
		logger: logger,
	}
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
