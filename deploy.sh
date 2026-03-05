#!/bin/bash

# 🚀 Script de Deploy Rápido - ProSync
# Este script automatiza o processo de deploy

set -e

echo "🚀 Iniciando deploy do ProSync..."

# Cores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Verificar se .env existe
if [ ! -f .env ]; then
    echo -e "${RED}❌ Arquivo .env não encontrado!${NC}"
    echo "Copie o arquivo .env.example para .env e configure:"
    echo "cp .env.example .env"
    exit 1
fi

# Perguntar qual ambiente
echo -e "${YELLOW}Qual ambiente deseja fazer deploy?${NC}"
echo "1) Desenvolvimento (docker-compose.yml)"
echo "2) Produção (docker-compose.prod.yml)"
read -p "Escolha (1 ou 2): " ambiente

if [ "$ambiente" = "1" ]; then
    COMPOSE_FILE="docker-compose.yml"
    echo -e "${GREEN}✅ Modo: Desenvolvimento${NC}"
elif [ "$ambiente" = "2" ]; then
    COMPOSE_FILE="docker-compose.prod.yml"
    echo -e "${GREEN}✅ Modo: Produção${NC}"
else
    echo -e "${RED}❌ Opção inválida!${NC}"
    exit 1
fi

# Parar containers existentes
echo -e "${YELLOW}🛑 Parando containers existentes...${NC}"
docker-compose -f $COMPOSE_FILE down || true

# Rebuild e start
echo -e "${YELLOW}🔨 Construindo imagens...${NC}"
docker-compose -f $COMPOSE_FILE build

echo -e "${YELLOW}🚀 Iniciando containers...${NC}"
docker-compose -f $COMPOSE_FILE up -d

# Aguardar alguns segundos
sleep 3

# Verificar status
echo -e "${YELLOW}📊 Verificando status dos containers...${NC}"
docker-compose -f $COMPOSE_FILE ps

# Testar health endpoint
echo -e "${YELLOW}🏥 Testando endpoint de saúde...${NC}"
sleep 2
if curl -f http://localhost:8000/api/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Aplicação está respondendo corretamente!${NC}"
else
    echo -e "${RED}⚠️  Aplicação pode não estar respondendo. Verifique os logs:${NC}"
    echo "docker-compose -f $COMPOSE_FILE logs -f"
fi

echo ""
echo -e "${GREEN}🎉 Deploy concluído!${NC}"
echo ""
echo "📝 Comandos úteis:"
echo "  Ver logs:      docker-compose -f $COMPOSE_FILE logs -f"
echo "  Parar:         docker-compose -f $COMPOSE_FILE down"
echo "  Reiniciar:     docker-compose -f $COMPOSE_FILE restart"
echo "  Status:        docker-compose -f $COMPOSE_FILE ps"
echo ""

if [ "$ambiente" = "1" ]; then
    echo "🌐 Acesse: http://localhost:8000"
else
    echo "🌐 Acesse seu domínio configurado no nginx"
    echo "   (verifique o arquivo nginx/nginx.conf)"
fi
