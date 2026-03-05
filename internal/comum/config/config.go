package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Tiny
	TinyBearerToken  string
	TinyIdListaPreco string
	TinyBaseURL      string

	// Bling (para uso futuro)
	BlingBaseURL string
	BlingAPIKey  string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Serviço
	IntervaloExecucaoMinutos int
	Ambiente                 string
	MaintenancePrice         float64
}

var configuracao *Config

// CarregarConfig carrega as configurações do arquivo .env
func CarregarConfig() (*Config, error) {
	// Tenta carregar o arquivo .env (ignora erro se não existir, usa variáveis do sistema)
	_ = godotenv.Load()

	config := &Config{
		// Tiny
		TinyBearerToken:  obterEnv("TINY_BEARER_TOKEN", ""),
		TinyIdListaPreco: obterEnv("TINY_ID_LISTA_PRECO", "43"),
		TinyBaseURL:      obterEnv("TINY_BASE_URL", "https://api.tiny.com.br/api2"),

		// Bling
		BlingBaseURL: obterEnv("BLING_BASE_URL", ""),
		BlingAPIKey:  obterEnv("BLING_API_KEY", ""),

		// Database
		DBHost:     obterEnv("DB_HOST", "localhost"),
		DBPort:     obterEnv("DB_PORT", "3306"),
		DBUser:     obterEnv("DB_USER", "root"),
		DBPassword: obterEnv("DB_PASSWORD", ""),
		DBName:     obterEnv("DB_NAME", "prosync"),

		// Serviço
		IntervaloExecucaoMinutos: obterEnvInt("INTERVALO_EXECUCAO_MINUTOS", 30),
		Ambiente:                 obterEnv("AMBIENTE", "development"),
		MaintenancePrice:         obterEnvFloat("MAINTENANCE_PRICE", 0),
	}

	// Valida configurações obrigatórias
	if config.TinyBearerToken == "" {
		return nil, fmt.Errorf("TINY_BEARER_TOKEN não configurado")
	}

	configuracao = config
	return config, nil
}

// ObterConfig retorna a configuração carregada
func ObterConfig() *Config {
	if configuracao == nil {
		panic("Configuração não foi carregada. Chame CarregarConfig() primeiro")
	}
	return configuracao
}

// obterEnv obtém uma variável de ambiente com valor padrão
func obterEnv(chave, valorPadrao string) string {
	if valor := os.Getenv(chave); valor != "" {
		return valor
	}
	return valorPadrao
}

// obterEnvInt obtém uma variável de ambiente inteira com valor padrão
func obterEnvInt(chave string, valorPadrao int) int {
	if valor := os.Getenv(chave); valor != "" {
		if intValor, err := strconv.Atoi(valor); err == nil {
			return intValor
		}
	}
	return valorPadrao
}

// obterEnvFloat obtém uma variável de ambiente float64 com valor padrão
func obterEnvFloat(chave string, valorPadrao float64) float64 {
	if valor := os.Getenv(chave); valor != "" {
		if floatValor, err := strconv.ParseFloat(valor, 64); err == nil {
			return floatValor
		}
	}
	return valorPadrao
}
