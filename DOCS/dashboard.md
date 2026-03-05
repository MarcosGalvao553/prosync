# 🌐 Dashboard Web - Guia de Uso

## Acessando o Dashboard

1. Execute o ProSync: `./prosync`
2. Aguarde a mensagem: `🌐 Servidor web iniciado em http://localhost:8080`
3. Abra seu navegador em: **http://localhost:8080**

## Interface do Dashboard

### 📊 Estatísticas Principais

No topo do dashboard, você verá 5 cards com métricas principais:
- **Total Requisições**: Número total de chamadas à API
- **Sucesso**: Requisições bem-sucedidas (status OK)
- **Erros**: Requisições com falha
- **Tempo Médio**: Média de duração das requisições em ms
- **Taxa de Sucesso**: Percentual de sucesso

### 🔍 Filtros Disponíveis

#### **Arquivo de Log**
Selecione qual arquivo de log você quer analisar. Os arquivos são nomeados por data (ex: `tiny_2026-03-03.json`)

#### **Serviço**
Filtre por serviço específico:
- `tiny` - Requisições para API do Tiny
- `sistema` - Logs do sistema
- `processador` - Logs do processador

#### **Operação**
Filtre por tipo de operação:
- `BuscarExcecoesListaPreco` - Busca de exceções de preço
- `BuscarDadosProduto` - Busca de dados do produto
- `BuscarEstoqueProduto` - Busca de estoque

#### **Status**
- `OK` - Requisições bem-sucedidas
- `Erro` - Requisições com erro

#### **ID do Produto**
Digite o ID de um produto específico para ver todas as requisições relacionadas a ele.

### 📈 Gráficos

#### **Distribuição de Status** (Gráfico de Pizza)
Mostra a proporção entre requisições bem-sucedidas e com erro.

#### **Tempo por Operação** (Gráfico de Barras)
Exibe o tempo médio de resposta para cada tipo de operação.

### 📝 Lista de Logs

Mostra as últimas 100 requisições filtradas. Para cada log:
- **Clique** para expandir e ver detalhes completos
- Veja **requisição**, **resposta**, **URL** e **erros** (se houver)

### Detalhes ao Clicar em um Log

Quando você clica em uma requisição, vê:

#### **Requisição**
```json
{
  "id": "967641344",
  "formato": "json",
  "token": "***OCULTO***"
}
```

#### **Resposta** (Exemplo - Produto)
```json
{
  "status": "OK",
  "status_processamento": "3",
  "id_produto": "967641344",
  "nome": "Funko Pop Jujutsu Kaisen..."
}
```

#### **Resposta** (Exemplo - Estoque)
```json
{
  "status": "OK",
  "id_produto": "967641344",
  "saldo": 24,
  "saldo_reservado": 14
}
```

## Exemplos de Uso

### 🔍 Analisar um produto específico

1. Digite o ID do produto no filtro "ID do Produto": `967641344`
2. Clique em "Aplicar Filtros"
3. Você verá:
   - Busca de exceção de preço para este produto
   - Dados completos do produto
   - Informações de estoque
   - Todos os detalhes de cada requisição

### 📊 Ver apenas erros

1. Selecione "Erro" no filtro Status
2. Clique em "Aplicar Filtros"
3. Analise quais produtos/operações estão falhando

### ⏱️ Monitorar performance

1. Observe o card "Tempo Médio"
2. Verifique o gráfico "Tempo por Operação"
3. Identifique operações lentas

### 🔄 Acompanhar em tempo real

- O dashboard atualiza automaticamente a cada 30 segundos
- Ou clique no botão "🔄 Atualizar" para atualizar manualmente

## API REST

O dashboard consome uma API REST que você também pode usar diretamente:

### Endpoints Disponíveis

#### `GET /api/health`
Verifica se o servidor está funcionando
```bash
curl http://localhost:8080/api/health
```

#### `GET /api/logs/arquivos`
Lista todos os arquivos de log disponíveis
```bash
curl http://localhost:8080/api/logs/arquivos
```

#### `GET /api/logs?[filtros]`
Retorna logs filtrados

**Parâmetros:**
- `arquivo` - Nome do arquivo (ex: `tiny_2026-03-03.json`)
- `servico` - Filtro por serviço
- `operacao` - Filtro por operação
- `status` - Filtro por status (OK ou Erro)
- `produto` - Filtro por ID do produto

**Exemplo:**
```bash
curl "http://localhost:8080/api/logs?servico=tiny&status=OK"
curl "http://localhost:8080/api/logs?produto=967641344"
```

## Dicas

✅ **Use filtros combinados**: Você pode combinar múltiplos filtros para análises específicas

✅ **Clique nos logs**: Sempre clique nos logs para ver os dados completos da requisição e resposta

✅ **Monitore a taxa de sucesso**: Se estiver abaixo de 95%, investigue os erros

✅ **Verifique tempos**: Se o tempo médio estiver alto, pode haver problemas de rede ou rate limiting

✅ **Analise por produto**: Use o filtro de ID do produto para debug específico
