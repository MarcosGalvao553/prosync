package repositories

import (
	"database/sql"
	"fmt"

	"prosync/internal/comum/models"
)

// PartnerRepository gerencia operações com partners
type PartnerRepository struct {
	db *sql.DB
}

// NovoPartnerRepository cria novo repositório de partners
func NovoPartnerRepository(db *sql.DB) *PartnerRepository {
	return &PartnerRepository{db: db}
}

// BuscarPorID busca um partner pelo ID
func (r *PartnerRepository) BuscarPorID(id int) (*models.Partner, error) {
	query := `
		SELECT id, razao_social, nome_fantasia, cnpj, cep, endereco, numero,
		       complemento, bairro, municipio, estado, nome_responsavel,
		       telefone_fixo, celular_whatsapp, email, desconto_negociado,
		       fee_adicional, repasse_comissao_mensal, url_redirect,
		       url_to_send_product, url_to_send_stock, created_at, updated_at
		FROM partners
		WHERE id = ?
	`

	var p models.Partner
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.RazaoSocial, &p.NomeFantasia, &p.CNPJ, &p.CEP, &p.Endereco,
		&p.Numero, &p.Complemento, &p.Bairro, &p.Municipio, &p.Estado,
		&p.NomeResponsavel, &p.TelefoneFixo, &p.CelularWhatsapp, &p.Email,
		&p.DescontoNegociado, &p.FeeAdicional, &p.RepasseComissaoMensal,
		&p.URLRedirect, &p.URLToSendProduct, &p.URLToSendStock,
		&p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("partner não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar partner: %w", err)
	}

	return &p, nil
}
