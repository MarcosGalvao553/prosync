# 🚀 Guia de Deploy - ProSync (Nginx Existente)

Este guia mostra como fazer o deploy do ProSync em um servidor que JÁ POSSUI um nginx rodando via Docker.

## 📋 Pré-requisitos

- Servidor Linux com Docker e Docker Compose
- Nginx já rodando (como no seu servidor nerdrop)
- Subdomínio configurado (ex: `prosync.nerdrop.com.br`)
- Acesso SSH ao servidor

## 🔧 Passo 1: Preparar o Servidor

Como você já tem Docker e nginx rodando, vamos apenas preparar o diretório do projeto.

## 🌐 Passo 2: Configurar DNS

No painel do seu provedor de domínio, crie um registro A:

```
Tipo: A
Nome: prosync (ou o subdomínio desejado)
Valor: IP_DO_SEU_SERVIDOR
TTL: 3600
```

Aguarde a propagação DNS (pode levar até 24h, mas geralmente leva minutos).

Teste com:
```bash
ping prosync.seudominio.com.br
```

## 📦 Passo 3: Fazer Upload do Projeto

### 3.1 No seu computador local:

```bash
# Criar arquivo compactado (excluindo arquivos desnecessários)
cd /home/galvao/Desktop/Freelas/nerdrop
tar -czf prosync.tar.gz prosync/ \
  --exclude='prosync/logs/*' \
  --exclude='prosync/.git' \ (nerdrop.com.br), crie um registro A:
 \
  --exclude='prosync/tmp/*'

# Enviar para o servidor
scp prosync.tar.gz root@IP_SERVIDOR:/workspace/
```

### 3.2 No servidor:

```bash
# Extrair projeto
cd /workspace
tar -xzf prosync.tar.gz
cd prosync

# Criar diretório de logs
mkdir -p logsnc.nerdrop
```

## ⚙️ Passo 4: Configurar Variáveis de Ambiente

```bash
# Copiar arquivo de exemplo
cp .env.example .env

# Editar com suas configurações (ajuste conforme seu banco):

```env
DB_HOST=banco
DB_PORT=3306
DB_DATABASE=nerdrop
DB_USERNAME=root
DB_PASSWORD=sua_senha_mysql

APP_PORT=8000
APP_ENV=production

DOMAIN=prosync.nerdrop.com.br
```

**Nota**: Use `DB_HOST=banNginx

### 5.1 Copiar configuração nginx (SEM SSL primeiro)

```bash
# Copiar arquivo de configuração do ProSync para o nginx existente
cp /workspace/prosync/nginx/prosync-sem-ssl.conf /workspace/drop/nginx/conf.d/prosync.conf

# Editar e ajustar o domínio se necessário
nano /workspace/drop/nginx/conf.d/prosync.conf
```

### 5.2 Recarregar nginx

```bash
# Recarregar configuração do nginx
docker exec drop-nginx-1 nginx -t  # Testar configuração
docker exec drop-nginx-1 nginx -s reload  # Recarregar
```

### 5.3 Gerar certificado SSL

```bash
# Entrar no container do nginx
docker exec -it drop-nginx-1 sh

# Dentro do container, gerar certificado
certbot certonly --webroot -w /var/www/html -d prosync.nerdrop.com.br

# Sair do container
exit
```

### 5.4 Ativar configuração com SSL

```bash
# Substituir pela configuração com SSL
cp /workspace/prosync/nginx/prosync.conf /workspace/drop/nginx/conf.d/prosync.conf

# Verificar e recarregar
docker exec drop-nginx-1 nginx -t
docker exec drop-nginx-1 nginx -s reloa
mv nginx.conf nginx-sem-ssl.conf
mv nginx-com-ssl.conf nginx.conf
naVoltar para o diretório do prosync
cd /workspace/prosync

# Construir e subir containeromínio

cd ..
docker-compose -f docker-compose.prod.yml up -d
```

## 🚀 Passo 6: Subir a Aplicação

```bash
# Construir e subir containers
docker-compose -f docker-compose.prod.yml up -d --build

# Verificar se está rodando
docker-compose -f docker-compose.prod.yml ps

# VeSem SSL**: `http://prosync.nerdrop.com.br`
- **Com SSL** (após configurar): `https://prosync.nerdrop.com.br`

O dashboard (monitor) estará disponível na página inicial!

### URLs disponíveis:
- Dashboard: `https://prosync.nerdrop.com.br/`
- Health: `https://prosync.nerdrop.com.br/api/health`
- Logs: `https://prosync.nerdrop.com.br/api/logs`
## 🎯 Passo 7: Acessar o Monitor

Abra seu navegador e acesse:

- **HTTP**: `http://prosync.seudominio.com.br` (será redirecionado para HTTPS)
- **HTTPS**: `https://prosync.seudominio.com.br`

O dashboard (monitor) estará disponível na página inicial!

## 🔍 Verificar Status

```bash
# Ver containers rodando
docker ps

# Ver logs em tempo real
docker-compose -f docker-compose.prod.yml logs -f
O certbot do nginx existente já deve ter renovação automática configurada. Caso não tenha:

```bash
# Adicionar cron job para renovação automática
crontab -e

# Adicionar esta linha:
0 0 * * 0 docker exec drop-nginx-1 certbot renew --quiet && docker exec drop-nginx-1 nginx -s reload
# Verificar saúde da aplicação
curl http://localhost:8000/api/health
```

## 🔄 Renovação Automática do SSL

``Ver todos os containers
docker ps

# Parar ProSync
cd /workspace/prosync
docker-compose -f docker-compose.prod.yml down

# Reiniciar ProSync
docker-compose -f docker-compose.prod.yml restart

# Reconstruir após mudanças no código
docker-compose -f docker-compose.prod.yml up -d --build

# Ver logs do ProSync
docker logs prosync -f

# Ver logs do 

Como você já tem outros serviços rodando, o firewall provavelmente já está configurado. Caso precise verificar:

```bash
# Verificar status do firewall

# Ver uso de recursos
docker statser-compose.prod.yml down

# Reiniciar aplicação
docker-compose -f docker-compose.prod.yml restart

# Reconstruir após mudanças no código
docker-compose -f docker-compose.prod.yml up -d --build

# Ver uso de recursos
docker stats

# Limpar containers antigos
docker system prune -a
```

## 🔥 Firewall (UFW)

```bash
# Configurar firewall
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
sudo ufw status
```

## 📊 Endpoints da API

Após o deploy, os seguintes endpoints estarão disponíveis:

- `GET /` - Dashboard (Monitor)
- `GET /api/health` - Verificação de saúde
- `GET /api/logs` - Logs da aplicação
- `GET /api/logs/estatisticas` - Estatísticas dos logs
- `GET /api/users` - Usuários
drop-nginx-1 nginx -t

# Ver logs do nginx
docker logs drop-nginx-1

# Verificar se a configuração do prosync está carregada
docker exec drop-nginx-1 cat /etc/nginx/conf.d/prosync.conf

```bash
docker-compose -f docker-compose.prod.yml logs prosync
```

### Erro de conexão com banco de dados

Verifique:
1. Credenciais no arquivo `.env`
2. Se o servidor do banco permite conexões externas
3. Firewall do servidor do banco de dados

### Nginx não responde

```bash
# Verificar configuração do nginx
docker exec prosync-nginx nginx -t

# Ver logs do nginnerdrop.com.br**

---

## 📝 Resumo Rápido dos Comandos

```bash
# No servidor (conectado via SSH)
cd /workspace/prosync

# 1. Configurar .env
cp .env.example .env
nano .env

# 2. Adicionar configuração nginx
cp nginx/prosync-sem-ssl.conf /workspace/drop/nginx/conf.d/prosync.conf
docker exec drop-nginx-1 nginx -s reload

# 3. Subir aplicação
docker-compose -f docker-compose.prod.yml up -d --build

# 4. Verificar logs
docker logs prosync -f

# 5. Testar
curl http://localhost:8000/api/health

# 6. Gerar SSL (após DNS propagado)
docker exec -it drop-nginx-1 certbot certonly --webroot -w /var/www/html -d prosync.nerdrop.com.br

# 7. Ativar SSL
cp nginx/prosync.conf /workspace/drop/nginx/conf.d/prosync.conf
docker exec drop-nginx-1 nginx -s reload
```
docker logs prosync-nginx
```

### SSL não funciona

1. Verifique se o DNS está propagado
2. Verifique se as portas 80 e 443 estão abertas
3. Tente gerar o certificado novamente

## 📝 Backup

```bash
# Criar backup dos logs
tar -czf backup-logs-$(date +%Y%m%d).tar.gz logs/

# Criar backup do banco (se local)
docker exec prosync-db mysqldump -u root -p nerdrop > backup-db-$(date +%Y%m%d).sql
```

## 🎉 Pronto!

Seu ProSync está rodando em produção! Acesse o monitor em:
**https://prosync.seudominio.com.br**
