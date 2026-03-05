package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// LogAPI representa um registro de log no banco de dados
type LogAPI struct {
	ID                 uint64
	CreatedAt          time.Time
	Servico            string
	Operacao           string
	Status             string
	RequestMethod      string
	RequestURL         string
	RequestHeaders     string
	RequestBody        string
	RequestSizeBytes   int
	ResponseStatusCode int
	ResponseHeaders    string
	ResponseBody       string
	ResponseSizeBytes  int
	DurationMs         float64
	ProdutoTinyID      string
	SKU                string
	IDProdutoBling     string
	UserID             uint64
	ErrorCode          string
	ErrorMessage       string
	Metadata           string
}

// SalvarLog insere um novo log no banco de dados
func SalvarLog(log *LogAPI) error {
	db := ObterConexao()

	query := `
		INSERT INTO logs_api (
			created_at, servico, operacao, status,
			request_method, request_url, request_headers, request_body, request_size_bytes,
			response_status_code, response_headers, response_body, response_size_bytes,
			duration_ms, produto_tiny_id, sku, idprodutobling, user_id,
			error_code, error_message, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query,
		log.CreatedAt,
		log.Servico,
		log.Operacao,
		log.Status,
		log.RequestMethod,
		log.RequestURL,
		log.RequestHeaders,
		log.RequestBody,
		log.RequestSizeBytes,
		log.ResponseStatusCode,
		log.ResponseHeaders,
		log.ResponseBody,
		log.ResponseSizeBytes,
		log.DurationMs,
		log.ProdutoTinyID,
		log.SKU,
		log.IDProdutoBling,
		log.UserID,
		log.ErrorCode,
		log.ErrorMessage,
		log.Metadata,
	)

	if err != nil {
		return fmt.Errorf("erro ao inserir log: %w", err)
	}

	lastID, _ := result.LastInsertId()
	log.ID = uint64(lastID)

	return nil
}

// BuscarLogs busca logs com filtros aplicados
func BuscarLogs(filtros map[string]string, limite int) ([]LogAPI, error) {
	db := ObterConexao()

	query := `
		SELECT 
			id, created_at, servico, operacao, status,
			request_method, request_url, request_headers, request_body, request_size_bytes,
			response_status_code, response_headers, response_body, response_size_bytes,
			duration_ms, produto_tiny_id, sku, idprodutobling, user_id,
			error_code, error_message, metadata
		FROM logs_api
		WHERE 1=1
	`

	args := []interface{}{}

	// Aplica filtros
	if servico, ok := filtros["servico"]; ok && servico != "" {
		query += " AND servico = ?"
		args = append(args, servico)
	}

	if operacao, ok := filtros["operacao"]; ok && operacao != "" {
		query += " AND operacao = ?"
		args = append(args, operacao)
	}

	if status, ok := filtros["status"]; ok && status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	if produtoID, ok := filtros["produto_tiny_id"]; ok && produtoID != "" {
		query += " AND produto_tiny_id = ?"
		args = append(args, produtoID)
	}

	if sku, ok := filtros["sku"]; ok && sku != "" {
		query += " AND sku = ?"
		args = append(args, sku)
	}

	if userID, ok := filtros["user_id"]; ok && userID != "" {
		query += " AND user_id = ?"
		args = append(args, userID)
	}

	if dataInicio, ok := filtros["data_inicio"]; ok && dataInicio != "" {
		query += " AND created_at >= ?"
		args = append(args, dataInicio)
	}

	if dataFim, ok := filtros["data_fim"]; ok && dataFim != "" {
		query += " AND created_at <= ?"
		args = append(args, dataFim)
	}

	// Ordena por data mais recente
	query += " ORDER BY created_at DESC"

	// Limite de resultados
	if limite > 0 {
		query += " LIMIT ?"
		args = append(args, limite)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar logs: %w", err)
	}
	defer rows.Close()

	logs := []LogAPI{}
	for rows.Next() {
		var log LogAPI
		err := rows.Scan(
			&log.ID,
			&log.CreatedAt,
			&log.Servico,
			&log.Operacao,
			&log.Status,
			&log.RequestMethod,
			&log.RequestURL,
			&log.RequestHeaders,
			&log.RequestBody,
			&log.RequestSizeBytes,
			&log.ResponseStatusCode,
			&log.ResponseHeaders,
			&log.ResponseBody,
			&log.ResponseSizeBytes,
			&log.DurationMs,
			&log.ProdutoTinyID,
			&log.SKU,
			&log.IDProdutoBling,
			&log.UserID,
			&log.ErrorCode,
			&log.ErrorMessage,
			&log.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// BuscarEstatisticas retorna estatísticas agregadas dos logs
func BuscarEstatisticas(filtros map[string]string) (map[string]interface{}, error) {
	db := ObterConexao()

	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'OK' THEN 1 ELSE 0 END) as sucesso,
			SUM(CASE WHEN status = 'Erro' THEN 1 ELSE 0 END) as erros,
			AVG(duration_ms) as tempo_medio,
			MIN(created_at) as primeira_requisicao,
			MAX(created_at) as ultima_requisicao
		FROM logs_api
		WHERE 1=1
	`

	args := []interface{}{}

	// Aplica filtros (mesma lógica da função BuscarLogs)
	if servico, ok := filtros["servico"]; ok && servico != "" {
		query += " AND servico = ?"
		args = append(args, servico)
	}

	if operacao, ok := filtros["operacao"]; ok && operacao != "" {
		query += " AND operacao = ?"
		args = append(args, operacao)
	}

	if status, ok := filtros["status"]; ok && status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	if produtoID, ok := filtros["produto_tiny_id"]; ok && produtoID != "" {
		query += " AND produto_tiny_id = ?"
		args = append(args, produtoID)
	}

	if sku, ok := filtros["sku"]; ok && sku != "" {
		query += " AND sku = ?"
		args = append(args, sku)
	}

	if userID, ok := filtros["user_id"]; ok && userID != "" {
		query += " AND user_id = ?"
		args = append(args, userID)
	}

	if dataInicio, ok := filtros["data_inicio"]; ok && dataInicio != "" {
		query += " AND created_at >= ?"
		args = append(args, dataInicio)
	}

	if dataFim, ok := filtros["data_fim"]; ok && dataFim != "" {
		query += " AND created_at <= ?"
		args = append(args, dataFim)
	}

	var total, sucesso, erros sql.NullInt64
	var tempoMedio sql.NullFloat64
	var primeiraReq, ultimaReq sql.NullTime

	err := db.QueryRow(query, args...).Scan(&total, &sucesso, &erros, &tempoMedio, &primeiraReq, &ultimaReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estatísticas: %w", err)
	}

	totalInt := int64(0)
	sucessoInt := int64(0)
	errosInt := int64(0)

	if total.Valid {
		totalInt = total.Int64
	}
	if sucesso.Valid {
		sucessoInt = sucesso.Int64
	}
	if erros.Valid {
		errosInt = erros.Int64
	}

	stats := map[string]interface{}{
		"total":        totalInt,
		"sucesso":      sucessoInt,
		"erros":        errosInt,
		"tempo_medio":  0.0,
		"taxa_sucesso": 0.0,
		"primeira_req": nil,
		"ultima_req":   nil,
	}

	if tempoMedio.Valid {
		stats["tempo_medio"] = tempoMedio.Float64
	}

	if totalInt > 0 {
		stats["taxa_sucesso"] = float64(sucessoInt) / float64(totalInt) * 100
	}

	if primeiraReq.Valid {
		stats["primeira_req"] = primeiraReq.Time
	}

	if ultimaReq.Valid {
		stats["ultima_req"] = ultimaReq.Time
	}

	return stats, nil
}

// BuscarTempoPorOperacao retorna tempo médio por tipo de operação
func BuscarTempoPorOperacao(filtros map[string]string) ([]map[string]interface{}, error) {
	db := ObterConexao()

	query := `
		SELECT operacao, AVG(duration_ms) as media
		FROM logs_api
		WHERE 1=1
	`

	args := []interface{}{}

	if servico, ok := filtros["servico"]; ok && servico != "" {
		query += " AND servico = ?"
		args = append(args, servico)
	}

	if operacao, ok := filtros["operacao"]; ok && operacao != "" {
		query += " AND operacao = ?"
		args = append(args, operacao)
	}

	if status, ok := filtros["status"]; ok && status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	if produtoID, ok := filtros["produto_tiny_id"]; ok && produtoID != "" {
		query += " AND produto_tiny_id = ?"
		args = append(args, produtoID)
	}

	if sku, ok := filtros["sku"]; ok && sku != "" {
		query += " AND sku = ?"
		args = append(args, sku)
	}

	if userID, ok := filtros["user_id"]; ok && userID != "" {
		query += " AND user_id = ?"
		args = append(args, userID)
	}

	if dataInicio, ok := filtros["data_inicio"]; ok && dataInicio != "" {
		query += " AND created_at >= ?"
		args = append(args, dataInicio)
	}

	if dataFim, ok := filtros["data_fim"]; ok && dataFim != "" {
		query += " AND created_at <= ?"
		args = append(args, dataFim)
	}

	query += " GROUP BY operacao ORDER BY media DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tempo por operação: %w", err)
	}
	defer rows.Close()

	resultados := []map[string]interface{}{}
	for rows.Next() {
		var operacao string
		var media float64
		if err := rows.Scan(&operacao, &media); err != nil {
			continue
		}
		resultados = append(resultados, map[string]interface{}{
			"operacao": operacao,
			"media":    media,
		})
	}

	return resultados, nil
}

// ConverterLogParaJSON converte um LogAPI para formato JSON legível
func ConverterLogParaJSON(log LogAPI) map[string]interface{} {
	resultado := map[string]interface{}{
		"id":          log.ID,
		"timestamp":   log.CreatedAt,
		"servico":     log.Servico,
		"operacao":    log.Operacao,
		"status":      log.Status,
		"status_code": log.ResponseStatusCode,
		"duracao":     fmt.Sprintf("%.2fms", log.DurationMs),
		"url":         log.RequestURL,
		"metodo_http": log.RequestMethod,
	}

	// Parse request headers
	if log.RequestHeaders != "" {
		var headers map[string]interface{}
		if err := json.Unmarshal([]byte(log.RequestHeaders), &headers); err == nil {
			resultado["request_headers"] = headers
		}
	}

	// Parse request body
	if log.RequestBody != "" {
		var body map[string]interface{}
		if err := json.Unmarshal([]byte(log.RequestBody), &body); err == nil {
			resultado["requisicao"] = body
		} else {
			resultado["requisicao"] = log.RequestBody
		}
	}

	// Parse response body
	if log.ResponseBody != "" {
		var body map[string]interface{}
		if err := json.Unmarshal([]byte(log.ResponseBody), &body); err == nil {
			resultado["resposta"] = body
		} else {
			resultado["resposta"] = log.ResponseBody
		}
	}

	// Parse metadata
	if log.Metadata != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(log.Metadata), &metadata); err == nil {
			resultado["metadata"] = metadata
		}
	}

	// Adiciona dados de produto se existirem
	if log.ProdutoTinyID != "" {
		resultado["produto_tiny_id"] = log.ProdutoTinyID
	}

	if log.SKU != "" {
		resultado["sku"] = log.SKU
	}

	if log.UserID > 0 {
		resultado["user_id"] = log.UserID
	}

	// Adiciona erro se existir
	if log.ErrorMessage != "" {
		resultado["erro"] = log.ErrorMessage
		if log.ErrorCode != "" {
			resultado["erro_codigo"] = log.ErrorCode
		}
	}

	return resultado
}

// BuscarUsuarios busca usuários do sistema para autocomplete
func BuscarUsuarios(termo string) ([]map[string]interface{}, error) {
	db := ObterConexao()

	query := `
		SELECT id, name, last_name, email, nickname
		FROM users
		WHERE 1=1
	`

	args := []interface{}{}

	// Se houver termo de busca, filtra por nome, email ou nickname
	if termo != "" {
		query += ` AND (
			name LIKE ? OR 
			last_name LIKE ? OR 
			email LIKE ? OR 
			nickname LIKE ?
		)`
		termoLike := "%" + termo + "%"
		args = append(args, termoLike, termoLike, termoLike, termoLike)
	}

	query += " ORDER BY name ASC LIMIT 50"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %w", err)
	}
	defer rows.Close()

	usuarios := []map[string]interface{}{}
	for rows.Next() {
		var id int64
		var name, lastName, email string
		var nickname sql.NullString

		err := rows.Scan(&id, &name, &lastName, &email, &nickname)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler usuário: %w", err)
		}

		usuario := map[string]interface{}{
			"id":    id,
			"name":  name + " " + lastName,
			"email": email,
		}

		if nickname.Valid {
			usuario["nickname"] = nickname.String
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}
