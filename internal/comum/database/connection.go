package database

import (
	"database/sql"
	"fmt"
	"time"

	"prosync/internal/comum/config"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InicializarConexao inicializa a conexão com o banco de dados
func InicializarConexao(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão com banco: %w", err)
	}

	// Configura pool de conexões
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Testa a conexão
	if err := db.Ping(); err != nil {
		return fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	return nil
}

// ObterConexao retorna a conexão ativa com o banco
func ObterConexao() *sql.DB {
	if db == nil {
		panic("Banco de dados não inicializado. Chame InicializarConexao() primeiro")
	}
	return db
}

// FecharConexao fecha a conexão com o banco
func FecharConexao() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
