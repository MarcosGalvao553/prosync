-- ==================================================
-- Migration: Adicionar campo user_id na tabela logs_api
-- Data: 2026-03-04
-- ==================================================

-- Adiciona a coluna user_id
ALTER TABLE logs_api 
ADD COLUMN user_id BIGINT UNSIGNED AFTER idprodutobling;

-- Adiciona índices para o novo campo
ALTER TABLE logs_api 
ADD INDEX idx_user_id (user_id),
ADD INDEX idx_user_servico (user_id, servico);

-- Verifica se foi adicionado corretamente
-- DESCRIBE logs_api;
