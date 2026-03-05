# ProSync

Microsserviço em Golang para sincronização de produtos entre Tiny e Bling.

## Estrutura do Projeto

```
prosync/
├── cmd/
│   └── main.go                 # Ponto de entrada (desacoplado)
├── internal/
│   ├── tiny/                   # Integração com Tiny
│   │   ├── dto/               # DTOs para requisição/resposta
│   │   │   ├── excecao_lista_preco.go
│   │   │   ├── produto.go
│   │   │   └── estoque.go
│   │   ├── entidade/          # Cliente HTTP e métodos
│   │   │   └── tiny_client.go
│   │   └── servico/           # Lógica de negócio
│   │       └── processador.go
│   ├── bling/                  # Integração com Bling (preparado)
│   │   ├── dto/
│   │   │   └── produto.go
│   │   └── entidade/
│   │       └── bling_client.go
│   └── comum/                  # Código compartilhado
│       ├── config/            # Gerenciamento de configurações
│       │   └── config.go
│       ├── logger/            # Sistema de logs
│       │   └── logger.go
│       └── servidor/          # Servidor HTTP
│           └── http.go
├── web/                        # Interface web
│   └── index.html             # Dashboard interativo
├── logs/                       # Logs gerados (criado automaticamente)
├── DOCS/                       # Documentação
├── .env                        # Configurações (não versionado)
├── .env.example               # Exemplo de configurações
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

## Configuração

1. Copie o arquivo `.env.example` para `.env`:
```bash
cp .env.example .env
```

2. Edite o arquivo `.env` com suas credenciais:
```env
TINY_BEARER_TOKEN=seu_token_aqui
TINY_ID_LISTA_PRECO=43
INTERVALO_EXECUCAO_MINUTOS=30
```

## Instalação

```bash
# Baixar dependências
go mod download

# Ou usar tidy para organizar
go mod tidy
```

## Execução

```bash
# Executar diretamente
go run cmd/main.go

# Ou compilar e executar
go build -o prosync cmd/main.go
./prosync
```

### 🌐 Dashboard Web

O serviço inclui um dashboard web interativo que inicia automaticamente em `http://localhost:8080`

**Recursos do Dashboard:**
- 📊 **Visualização gráfica** de todos os logs
- 🔍 **Filtros avançados** por serviço, operação, status, produto
- 📈 **Gráficos em tempo real** de performance e status
- 📝 **Lista detalhada** de todas as requisições
- 🔄 **Atualização automática** a cada 30 segundos
- 💾 **Seleção de arquivos** de log por data

**Como usar:**
1. Execute o ProSync: `./prosync`
2. Abra o navegador em: `http://localhost:8080`
3. Selecione o arquivo de log desejado
4. Use os filtros para análise específica
5. Clique em qualquer log para ver detalhes completos

## Funcionalidades

### 🌐 Dashboard Web Interativo

- ✅ Interface web moderna e responsiva
- ✅ Servidor HTTP integrado (porta 8080)
- ✅ Filtros em tempo real por:
  - Arquivo de log (por data)
  - Serviço (tiny, bling, sistema)
  - Operação (BuscarExcecoesListaPreco, BuscarDadosProduto, etc)
  - Status (OK, Erro)
  - ID do Produto
- ✅ Visualização detalhada de requisições e respostas
- ✅ Gráficos de performance e distribuição
- ✅ Atualização automática

### Integração com Tiny

- ✅ Buscar exceções de lista de preços (com paginação automática)
- ✅ Buscar dados completos de produtos
- ✅ Buscar estoque de produtos
- ✅ Rate limiting automático (1 req/s)
- ✅ Retry automático em caso de rate limit excedido
- ✅ Logs detalhados de todas as requisições

### Fluxo de Processamento

O serviço executa o seguinte fluxo automaticamente:

1. **Busca exceções de lista de preços** - Itera por todas as páginas
2. **Para cada produto encontrado:**
   - Busca dados completos do produto
   - Busca informações de estoque
   - Armazena tudo em uma coleção unificada
3. **Gera estatísticas** do processamento
4. **Aguarda** o intervalo configurado e reinicia

### Rate Limiting

- **Limite normal:** Máximo de 1 requisição por segundo
- **Rate limit excedido:** Aguarda 1 minuto automaticamente e retenta
- Código de erro 6 da API = Rate limit atingido

### Sistema de Logs

Os logs são salvos na pasta `logs/` em dois formatos:

1. **JSON** (`tiny_2026-03-03.json`): Estruturado para parsing
2. **Texto** (`tiny_2026-03-03.log`): Legível para humanos

Cada log contém:
- Timestamp
- URL e método HTTP
- Dados da requisição (token oculto)
- Resposta completa
- Status code
- Duração da requisição
- Erros (se houver)

### Configurações

Todas as configurações são feitas via variáveis de ambiente:

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `TINY_BEARER_TOKEN` | Token de autenticação Tiny | - |
| `TINY_ID_LISTA_PRECO` | ID da lista de preços | 43 |
| `TINY_BASE_URL` | URL base da API Tiny | https://api.tiny.com.br/api2 |
| `INTERVALO_EXECUCAO_MINUTOS` | Intervalo entre execuções | 30 |
| `AMBIENTE` | Ambiente (development/production) | development |

## Desenvolvimento

### Adicionar nova integração

1. Crie a pasta em `internal/nome_servico/`
2. Adicione DTOs em `internal/nome_servico/dto/`
3. Crie a entidade cliente em `internal/nome_servico/entidade/`
4. Implemente os métodos HTTP seguindo o padrão do TinyClient

### Logs

O logger registra automaticamente:
- Todas as chamadas HTTP (request/response)
- Informações importantes do processamento
- Erros e exceções

Os arquivos são separados por serviço e por dia.

## Próximos Passos

- [ ] Implementar integração com Bling
- [ ] Adicionar persistência no MySQL
- [ ] Implementar sincronização entre Tiny e Bling
- [ ] Adicionar métricas de performance
- [ ] Implementar circuit breaker
- [ ] Adicionar testes unitários
- [ ] Dockerizar a aplicação
- [ ] API REST para consulta de dados processados

## Arquitetura

### Camadas

1. **CMD** - Ponto de entrada, mantém código mínimo
2. **Serviço** - Orquestra a lógica de negócio
3. **Entidade** - Clientes HTTP e comunicação com APIs
4. **DTO** - Estruturas de dados de entrada/saída
5. **Comum** - Recursos compartilhados (config, logger)

### Princípios

- **Desacoplamento**: Main limpo, lógica no serviço
- **Single Responsibility**: Cada camada tem uma responsabilidade
- **Rate Limiting**: Respeita limites da API automaticamente
- **Retry**: Tratamento automático de erros temporários
- **Logging**: Todas as operações são logadas

Estrutura que agrega todas as informações de um produto:

```go
type ProdutoCompleto struct {
    Excecao dto.ProdutoExcecaoListaPrecoTiny  // Dados da exceção de preço
    Produto *dto.ProdutoTiny                   // Dados completos do produto
    Estoque *dto.EstoqueTiny                   // Informações de estoque
}
```

### ProdutoExcecaoListaPrecoTiny

Dados da exceção de lista de preços:
- `ID`: Identificador do registro
- `IdProduto`: Identificador do produto no Tiny
- `Preco`: Preço da exceção

### ProdutoTiny

Dados essenciais do produto:
- `ID`, `Nome`, `Codigo`
- `Preco`, `PrecoCusto`
- `Marca`, `Categoria`, `Situacao`, `GTIN`

### EstoqueTiny

Informações de estoque:
- `IDProduto`, `Nome`
- `Saldo`, `SaldoReservado`, `SaldoDisponivel`

Estas estruturas contêm apenas os dados essenciais extraídos das APIs para facilitar o
- `IdProduto`: Identificador do produto no Tiny
- `Preco`: Preço da exceção

Este objeto é gerado a partir da resposta completa da API e contém apenas os dados essenciais para processamento.
