package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"prosync/internal/comum/models"
)

// BlingConfigurationRepository gerencia operações com configurações do Bling
type BlingConfigurationRepository struct {
	db *sql.DB
}

// NovoBlingConfigurationRepository cria novo repositório
func NovoBlingConfigurationRepository(db *sql.DB) *BlingConfigurationRepository {
	return &BlingConfigurationRepository{db: db}
}

// BuscarPorUserID busca configuração do Bling por ID do usuário
func (r *BlingConfigurationRepository) BuscarPorUserID(userID int) (*models.BlingConfiguration, error) {
	query := `
		SELECT id, client_id, secret_key, url_callback, postcode, access_token, 
		       refresh_token, token_validate, code, user_id, created_at, updated_at
		FROM bling_configurations 
		WHERE user_id = ?
		LIMIT 1
	`

	var config models.BlingConfiguration
	err := r.db.QueryRow(query, userID).Scan(
		&config.ID, &config.ClientID, &config.SecretKey, &config.URLCallback,
		&config.Postcode, &config.AccessToken, &config.RefreshToken,
		&config.TokenValidate, &config.Code, &config.UserID,
		&config.CreatedAt, &config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Retorna nil sem erro quando não encontra
		}
		return nil, fmt.Errorf("erro ao buscar configuração do Bling: %w", err)
	}

	return &config, nil
}

// AtualizarTokens atualiza access_token, refresh_token e token_validate
func (r *BlingConfigurationRepository) AtualizarTokens(id int, accessToken, refreshToken string, tokenValidate time.Time) error {
	query := `
		UPDATE bling_configurations 
		SET access_token = ?, refresh_token = ?, token_validate = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.Exec(query, accessToken, refreshToken, tokenValidate, id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar tokens do Bling: %w", err)
	}

	return nil
}

// TokenEstaValido verifica se o token ainda está válido (com margem de 1 hora)
func (r *BlingConfigurationRepository) TokenEstaValido(config *models.BlingConfiguration) bool {
	if !config.TokenValidate.Valid {
		return false
	}

	// Token deve ter pelo menos 1 hora de validade restante
	return time.Now().Add(1 * time.Hour).Before(config.TokenValidate.Time)
}
