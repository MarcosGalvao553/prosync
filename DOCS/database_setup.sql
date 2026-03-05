-- ==================================================
-- Script de Criação da Tabela de Logs do ProSync
-- ==================================================

CREATE TABLE logs_api (
    -- Identificação
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    
    -- Categorização
    servico VARCHAR(50) NOT NULL,                    -- 'tiny', 'bling', 'sistema', 'processador'
    operacao VARCHAR(100) NOT NULL,                  -- 'BuscarExcecoesListaPreco', 'BuscarDadosProduto', etc
    status ENUM('OK', 'Erro') NOT NULL,
    
    -- Requisição (dados completos, sem ocultação)
    request_method VARCHAR(10),                      -- 'GET', 'POST', etc
    request_url VARCHAR(500),                        -- URL completa
    request_headers JSON,                            -- Headers completos (incluindo token real)
    request_body LONGTEXT,                           -- Body completo da requisição
    request_size_bytes INT UNSIGNED,                 -- Tamanho da requisição em bytes
    
    -- Resposta (dados completos)
    response_status_code INT,                        -- 200, 404, 500, etc
    response_headers JSON,                           -- Headers da resposta
    response_body LONGTEXT,                          -- Response completa (pode ser muito grande)
    response_size_bytes INT UNSIGNED,                -- Tamanho da resposta em bytes
    
    -- Performance
    duration_ms DECIMAL(10,2),                       -- Duração em milissegundos (precisão de 2 casas)
    
    -- Dados de negócio (facilita filtros)
    produto_tiny_id VARCHAR(50),                     -- ID do produto no Tiny
    sku VARCHAR(100),                                -- SKU do produto
    idprodutobling VARCHAR(50),                      -- ID do produto no Bling (opcional por enquanto)
    user_id BIGINT UNSIGNED,                         -- ID do product_user (para filtrar requests do Bling por usuário)
    
    -- Erro (quando houver)
    error_code VARCHAR(50),                          -- Código de erro da API
    error_message TEXT,                              -- Mensagem de erro
    
    -- Contexto adicional (flexível)
    metadata JSON,                                   -- Qualquer dado extra (paginação, flags, etc)
    
    -- Índices para otimizar queries do dashboard
    INDEX idx_created_at (created_at),
    INDEX idx_servico (servico),
    INDEX idx_operacao (operacao),
    INDEX idx_status (status),
    INDEX idx_produto_tiny_id (produto_tiny_id),
    INDEX idx_sku (sku),
    INDEX idx_idprodutobling (idprodutobling),
    INDEX idx_user_id (user_id),
    INDEX idx_servico_status (servico, status),
    INDEX idx_operacao_status (operacao, status),
    INDEX idx_created_servico (created_at, servico),
    INDEX idx_user_servico (user_id, servico),
    
    -- Índice composto para queries complexas do dashboard
    INDEX idx_dashboard (created_at DESC, servico, operacao, status, produto_tiny_id)
    
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==================================================
-- Queries úteis para manutenção
-- ==================================================

-- Ver últimos 10 logs
-- SELECT id, created_at, servico, operacao, status, duration_ms, produto_tiny_id 
-- FROM logs_api 
-- ORDER BY created_at DESC 
-- LIMIT 10;

-- Estatísticas do dia
-- SELECT 
--     servico,
--     COUNT(*) as total,
--     SUM(CASE WHEN status = 'OK' THEN 1 ELSE 0 END) as sucesso,
--     SUM(CASE WHEN status = 'Erro' THEN 1 ELSE 0 END) as erros,
--     AVG(duration_ms) as tempo_medio_ms
-- FROM logs_api
-- WHERE DATE(created_at) = CURDATE()
-- GROUP BY servico;

-- Limpar logs antigos (exemplo: mais de 30 dias)
-- DELETE FROM logs_api WHERE created_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
