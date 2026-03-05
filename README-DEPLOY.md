# 🚀 ProSync - Guia Rápido de Deploy

## Para acessar o MONITOR do ProSync:

Acesse pelo navegador: **https://prosync.nerdrop.com.br**

O monitor é o dashboard na página inicial que mostra:
- Estatísticas de sincronização
- Logs em tempo real
- Gráficos de desempenho
- Status do sistema

---

## 📦 Como fazer o deploy:

### 1️⃣ No seu computador (enviar para servidor):

```bash
cd /home/galvao/Desktop/Freelas/nerdrop
tar -czf prosync.tar.gz prosync/ --exclude='prosync/logs/*' --exclude='prosync/.git'
scp prosync.tar.gz root@SEU_IP:/workspace/
```

### 2️⃣ No servidor (via SSH):

```bash
# Extrair projeto
cd /workspace
tar -xzf prosync.tar.gz
cd prosync

# Configurar variáveis
cp .env.example .env
nano .env  # Edite com suas credenciais do banco

# Executar deploy automatizado
chmod +x deploy-server.sh
./deploy-server.sh
```

### 3️⃣ Configurar DNS:

No painel do seu domínio:
- **Tipo**: A
- **Nome**: prosync
- **Valor**: IP do servidor

### 4️⃣ Gerar SSL (após DNS propagar):

```bash
# Gerar certificado
docker exec -it drop-nginx-1 certbot certonly --webroot -w /var/www/html -d prosync.nerdrop.com.br

# Ativar SSL
cp /workspace/prosync/nginx/prosync.conf /workspace/drop/nginx/conf.d/prosync.conf
docker exec drop-nginx-1 nginx -s reload
```

### 5️⃣ Acessar:

**https://prosync.nerdrop.com.br**

---

## 📊 Endpoints da API:

- `GET /` - Dashboard (Monitor)
- `GET /api/health` - Status da aplicação
- `GET /api/logs` - Logs de sincronização
- `GET /api/logs/estatisticas` - Estatísticas
- `GET /api/users` - Usuários

---

## 🔧 Comandos úteis:

```bash
# Ver logs
docker logs prosync -f

# Reiniciar
cd /workspace/prosync
docker-compose -f docker-compose.prod.yml restart

# Parar
docker-compose -f docker-compose.prod.yml down

# Rebuild
docker-compose -f docker-compose.prod.yml up -d --build

# Status
docker ps | grep prosync
```

---

## ⚙️ Variáveis importantes no .env:

```env
DB_HOST=banco              # Nome do container MySQL
DB_PORT=3306
DB_DATABASE=nerdrop
DB_USERNAME=root
DB_PASSWORD=sua_senha
```

---

## 🆘 Troubleshooting:

**Container não inicia?**
```bash
docker logs prosync
```

**Erro de conexão com banco?**
- Verifique credenciais no `.env`
- Teste: `docker exec prosync ping banco`

**Nginx não responde?**
```bash
docker logs drop-nginx-1
docker exec drop-nginx-1 nginx -t
```

---

## 📁 Estrutura do projeto:

```
prosync/
├── docker-compose.prod.yml  # Configuração Docker para produção
├── .env                     # Variáveis de ambiente (criar)
├── nginx/
│   ├── prosync.conf         # Config nginx com SSL
│   └── prosync-sem-ssl.conf # Config nginx sem SSL
├── deploy-server.sh         # Script automático de deploy
└── DEPLOY.md               # Guia completo de deploy
```

---

**Pronto! Seu monitor estará acessível em https://prosync.nerdrop.com.br** 🎉
