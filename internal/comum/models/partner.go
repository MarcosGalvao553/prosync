package models

import "time"

// Partner representa um parceiro/cliente do sistema
type Partner struct {
	ID                      int       `db:"id" json:"id"`
	RazaoSocial             string    `db:"razao_social" json:"razao_social"`
	NomeFantasia            *string   `db:"nome_fantasia" json:"nome_fantasia,omitempty"`
	CNPJ                    string    `db:"cnpj" json:"cnpj"`
	CEP                     string    `db:"cep" json:"cep"`
	Endereco                string    `db:"endereco" json:"endereco"`
	Numero                  *string   `db:"numero" json:"numero,omitempty"`
	Complemento             *string   `db:"complemento" json:"complemento,omitempty"`
	Bairro                  string    `db:"bairro" json:"bairro"`
	Municipio               string    `db:"municipio" json:"municipio"`
	Estado                  string    `db:"estado" json:"estado"`
	NomeResponsavel         string    `db:"nome_responsavel" json:"nome_responsavel"`
	TelefoneFixo            *string   `db:"telefone_fixo" json:"telefone_fixo,omitempty"`
	CelularWhatsapp         string    `db:"celular_whatsapp" json:"celular_whatsapp"`
	Email                   string    `db:"email" json:"email"`
	DescontoNegociado       float64   `db:"desconto_negociado" json:"desconto_negociado"`
	FeeAdicional            float64   `db:"fee_adicional" json:"fee_adicional"`
	RepasseComissaoMensal   float64   `db:"repasse_comissao_mensal" json:"repasse_comissao_mensal"`
	URLRedirect             *string   `db:"url_redirect" json:"url_redirect,omitempty"`
	URLToSendProduct        *string   `db:"url_to_send_product" json:"url_to_send_product,omitempty"`
	URLToSendStock          *string   `db:"url_to_send_stock" json:"url_to_send_stock,omitempty"`
	CreatedAt               time.Time `db:"created_at" json:"created_at"`
	UpdatedAt               time.Time `db:"updated_at" json:"updated_at"`
}

// TableName retorna o nome da tabela no banco de dados
func (Partner) TableName() string {
	return "partners"
}
