# 🐳 Guia Docker - ProSync

## 📋 Arquivos Docker

- **Dockerfile** - Produção (multi-stage build otimizado)
- **Dockerfile.dev** - Desenvolvimento (com hot reload)
- **docker-compose.yml** - Desenvolvimento local
- **docker-compose.prod.yml** - Produção no servidor

## 🔧 Desenvolvimento Local

```bash
# Subir em modo desenvolvimento (com hot reload)
docker-compose up

# Rebuild
docker-compose up --build

# Ver logs
docker-compose logs -f

# Parar
docker-compose down
```

## 🚀 Produção

```bash
# No servidor
docker-compose -f docker-compose.prod.yml up -d --build

# Ver logs
docker-compose -f docker-compose.prod.yml logs -f

# Reiniciar
docker-compose -f docker-compose.prod.yml restart

# Parar
docker-compose -f docker-compose.prod.yml down
```

## 🔑 Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto:

```env
# Banco de Dados
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=nerdrop
DB_USERNAME=root
DB_PASSWORD=sua_senha

# API Tiny
TINY_BEARER_TOKEN=seu_token
TINY_ID_LISTA_PRECO=43
TINY_BASE_URL=https://api.tiny.com.br/api2

# Configurações
INTERVALO_EXECUCAO_MINUTOS=30
AMBIENTE=development
```

## 📦 Build Manual

```bash
# Build da imagem
docker build -t prosync:latest -f Dockerfile .

# Rodar container
docker run -p 8000:8000 \
  -e DB_HOST=localhost \
  -e DB_PORT=3306 \
  -e DB_DATABASE=nerdrop \
  -e DB_USERNAME=root \
  -e DB_PASSWORD=senha \
  -v $(pwd)/logs:/app/logs \
  prosync:latest
```

## 🔍 Troubleshooting

### Container não inicia

```bash
# Ver logs detalhados
docker logs prosync

# Entrar no container
docker exec -it prosync sh
```

### Erro de build

```bash
# Limpar cache e rebuild
docker-compose build --no-cache

# Remover volumes órfãos
docker-compose down -v
```

### Porta 8000 já em uso

```bash
# Ver o que está usando a porta
lsof -i :8000

# Mudar porta no docker-compose.yml:
ports:
  - "8080:8000"  # Acesse em localhost:8080
```

## 📊 Otimizações do Dockerfile de Produção

- ✅ Multi-stage build (imagem final ~20MB)
- ✅ Build estático sem CGO
- ✅ Apenas binário compilado na imagem final
- ✅ Certificados SSL incluídos
- ✅ Pasta web (dashboard) copiada
- ✅ Diretório de logs criado
- ✅ Usa Alpine Linux (mínimo)

## 🆚 Diferenças Dev vs Prod

| Recurso | Desenvolvimento | Produção |
|---------|----------------|----------|
| Dockerfile | Dockerfile.dev | Dockerfile |
| Tamanho imagem | ~500MB | ~30MB |
| Hot reload | ✅ Sim (air) | ❌ Não |
| Build | On-demand | Multi-stage |
| Volumes | Todo código | Apenas logs |
| Network | Bridge | drop_default |

## 🌐 Acesso

- **Desenvolvimento**: http://localhost:8000
- **Produção**: https://prosync.nerdrop.com.br

O monitor/dashboard está disponível na raiz (/) após subir o container.
