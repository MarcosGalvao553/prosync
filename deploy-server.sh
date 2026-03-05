#!/bin/bash

# 🚀 Script de Deploy Rápido - ProSync (Servidor com Nginx Existente)
# Execute este script NO SERVIDOR após fazer upload do projeto

set -e

echo "🚀 Deploy ProSync - Usando nginx existente"
echo "=========================================="

# Verificar se está no diretório correto
if [ ! -f "docker-compose.prod.yml" ]; then
    echo "❌ Erro: Execute este script no diretório /workspace/prosync"
    exit 1
fi

# Verificar se .env existe
if [ ! -f .env ]; then
    echo "⚠️  Arquivo .env não encontrado!"
    echo "📝 Criando .env a partir do .env.example..."
    cp .env.example .env
    echo "✅ Arquivo .env criado. EDITE-O antes de continuar!"
    echo "   nano .env"
    exit 1
fi

# Perguntar se deseja continuar
read -p "⚠️  Certifique-se de que o .env está configurado. Continuar? (s/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Ss]$ ]]; then
    echo "❌ Deploy cancelado."
    exit 1
fi

# Criar diretório de logs se não existir
mkdir -p logs

echo ""
echo "📋 Passo 1: Configurando nginx..."
echo "Copiando configuração nginx (sem SSL)..."

if [ ! -d "/workspace/drop/nginx/conf.d" ]; then
    echo "❌ Erro: Diretório /workspace/drop/nginx/conf.d não encontrado!"
    echo "   Verifique se o nginx está configurado corretamente."
    exit 1
fi

cp nginx/prosync-sem-ssl.conf /workspace/drop/nginx/conf.d/prosync.conf
echo "✅ Configuração copiada para /workspace/drop/nginx/conf.d/prosync.conf"

# Testar configuração nginx
echo ""
echo "🔍 Testando configuração nginx..."
if docker exec drop-nginx-1 nginx -t; then
    echo "✅ Configuração nginx OK"
    echo "🔄 Recarregando nginx..."
    docker exec drop-nginx-1 nginx -s reload
    echo "✅ Nginx recarregado"
else
    echo "❌ Erro na configuração nginx!"
    exit 1
fi

# Construir e subir container prosync
echo ""
echo "🔨 Passo 2: Construindo e subindo ProSync..."
docker-compose -f docker-compose.prod.yml down 2>/dev/null || true
docker-compose -f docker-compose.prod.yml up -d --build

# Aguardar container iniciar
echo ""
echo "⏳ Aguardando aplicação iniciar..."
sleep 5

# Verificar se está rodando
echo ""
echo "📊 Status dos containers:"
docker-compose -f docker-compose.prod.yml ps

# Testar health endpoint
echo ""
echo "🏥 Testando endpoint de saúde..."
sleep 3

if curl -f http://localhost:8000/api/health > /dev/null 2>&1; then
    echo "✅ Aplicação está respondendo corretamente!"
else
    echo "⚠️  Aplicação pode não estar respondendo."
    echo "   Verifique os logs: docker logs prosync -f"
fi

echo ""
echo "🎉 Deploy concluído!"
echo "=========================================="
echo ""
echo "📝 Próximos passos:"
echo ""
echo "1. Configure o DNS (se ainda não configurou):"
echo "   - Tipo: A"
echo "   - Nome: prosync"
echo "   - Valor: IP do servidor"
echo ""
echo "2. Teste o acesso:"
echo "   http://prosync.nerdrop.com.br"
echo ""
echo "3. Configure SSL (após DNS propagado):"
echo "   docker exec -it drop-nginx-1 certbot certonly --webroot -w /var/www/html -d prosync.nerdrop.com.br"
echo ""
echo "4. Ative configuração com SSL:"
echo "   cp /workspace/prosync/nginx/prosync.conf /workspace/drop/nginx/conf.d/prosync.conf"
echo "   docker exec drop-nginx-1 nginx -s reload"
echo ""
echo "📚 Comandos úteis:"
echo "  Ver logs:     docker logs prosync -f"
echo "  Parar:        docker-compose -f docker-compose.prod.yml down"
echo "  Reiniciar:    docker-compose -f docker-compose.prod.yml restart"
echo ""
