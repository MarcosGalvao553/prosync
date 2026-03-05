# Models - Modelos de Banco de Dados

Modelos Go convertidos dos modelos PHP (Laravel/Eloquent) para uso com banco de dados MySQL/MariaDB.

## Estrutura

Todos os modelos estão no pacote `models` e representam as tabelas do banco de dados.

## Modelos Disponíveis

### 1. Product (`product.go`)
Representa produtos no sistema
- **Tabela**: `products`
- **Campos principais**: ID, Nome, Descrição, Preço, SKU, Estoque, etc.
- **Timestamps**: created_at, updated_at

### 2. Category (`category.go`)
Categorias de produtos
- **Tabela**: `categories`
- **Campos principais**: ID, Nome, Código, Imagem, Ranges de preço

### 3. ProductImage (`product_image.go`)
Imagens associadas a produtos
- **Tabela**: `product_images`
- **Campos principais**: ID, ProductID, ImageSrc, ImageType

### 4. ProductUser (`product_user.go`)
Relação entre produtos e usuários (preços customizados, etc)
- **Tabela**: `product_users`
- **Campos principais**: UserID, ProductID, TinyProductID, BlingProductID, Price

### 5. BlingConfiguration (`bling_configuration.go`)
Configurações do Bling por usuário
- **Tabela**: `bling_configurations`
- **Campos principais**: ClientID, SecretKey, AccessToken, RefreshToken

### 6. PreSaleProduct (`pre_sale_product.go`)
Produtos em pré-venda
- **Tabela**: `pre_sale_products`
- **Campos principais**: ProductID, EndDate, Active

### 7. TinyOrder (`tiny_order.go`)
Pedidos do Tiny
- **Tabela**: `tiny_orders`
- **Campos principais**: OrderTinyID, ShippingOrderID

### 8. SystemConfig (`system_config.go`)
Configurações gerais do sistema
- **Tabela**: `system_configs`
- **Campos principais**: Name, Description

### 9. SystemConfigParam (`system_config_param.go`)
Parâmetros de configuração
- **Tabela**: `system_config_params`
- **Campos principais**: Name, Code, SystemConfigID, ShowToUser

### 10. SystemConfigParamValue (`system_config_param_value.go`)
Valores dos parâmetros de configuração
- **Tabela**: `system_config_param_values`
- **Campos principais**: Name, Value, UserID, SystemConfigParamID

### 11. Config (`config.go`)
Configurações simples (chave-valor)
- **Tabela**: `configs`
- **Campos principais**: Code, Description, Value

## Uso

### Exemplo de consulta básica:

```go
package main

import (
    "database/sql"
    "prosync/internal/comum/database"
    "prosync/internal/comum/models"
)

func BuscarProduto(id int) (*models.Product, error) {
    db := database.ObterConexao()
    
    query := "SELECT * FROM products WHERE id = ?"
    
    var produto models.Product
    err := db.QueryRow(query, id).Scan(
        &produto.ID,
        &produto.Name,
        &produto.Description,
        &produto.Price,
        &produto.CostPrice,
        &produto.CategoryID,
        &produto.IsEnabled,
        &produto.IsPreSale,
        &produto.SaleCount,
        &produto.ReviewCount,
        &produto.SKU,
        &produto.Stock,
        &produto.Observation,
        &produto.Weight,
        &produto.Height,
        &produto.Width,
        &produto.Length,
        &produto.ProductTiny,
        &produto.NCM,
        &produto.EAN,
        &produto.Marca,
        &produto.CEST,
        &produto.StopStock,
        &produto.PromotionID,
        &produto.OriginalPrice,
        &produto.CreatedAt,
        &produto.UpdatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    return &produto, nil
}
```

### Exemplo de inserção:

```go
func CriarCategoria(nome string) error {
    db := database.ObterConexao()
    
    query := "INSERT INTO categories (name) VALUES (?)"
    _, err := db.Exec(query, nome)
    
    return err
}
```

## Tipos SQL Nullable

Os modelos utilizam tipos `sql.Null*` para campos que podem ser NULL no banco:
- `sql.NullString` - para strings nullable
- `sql.NullInt64` - para inteiros nullable
- `sql.NullFloat64` - para floats nullable
- `sql.NullBool` - para booleans nullable
- `sql.NullTime` - para timestamps nullable

### Como usar:

```go
produto.SKU.Valid // true se o valor não é NULL
produto.SKU.String // valor da string

// Para setar um valor:
produto.SKU = sql.NullString{String: "ABC123", Valid: true}

// Para setar NULL:
produto.SKU = sql.NullString{Valid: false}
```

## Tags

Os modelos incluem tags para:
- **`db`**: nome da coluna no banco de dados
- **`json`**: nome do campo ao serializar para JSON

## Próximos Passos

Para facilitar o uso, considere criar:
1. **Repositórios** - funções específicas para cada modelo (CRUD)
2. **Migrations** - scripts para criar/atualizar tabelas
3. **Testes** - testes unitários para validar os modelos
