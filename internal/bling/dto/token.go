package dto

// TokenRefreshRequest representa a requisição de refresh do token
type TokenRefreshRequest struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

// TokenRefreshResponse representa a resposta do refresh do token
type TokenRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // Segundos até expirar
	TokenType    string `json:"token_type"`
}

// ErroAPIBling representa um erro retornado pela API do Bling
type ErroAPIBling struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Type        string            `json:"type"`
	Message     string            `json:"message"`
	Description string            `json:"description,omitempty"`
	Fields      map[string]string `json:"fields,omitempty"`
}
