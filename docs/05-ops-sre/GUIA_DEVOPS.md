> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 🚀 Guia DevOps

**Versão:** 1.0  
**Data:** 14/11/2025  
**Status:** Guia Prático

---

## 📋 Índice

1. [Infraestrutura](#infraestrutura)
2. [NGINX](#nginx)
3. [CI/CD](#cicd)
4. [Monitoramento](#monitoramento)
5. [Backup & Recovery](#backup--recovery)
6. [Troubleshooting](#troubleshooting)

---

## 🏗️ Infraestrutura

### VPS Requirements

```
Mínimo:
- 2 vCPU
- 4 GB RAM
- 50 GB SSD
- Ubuntu 22.04 LTS

Recomendado (Produção):
- 4 vCPU
- 8+ GB RAM
- 100+ GB SSD
- Ubuntu 22.04 LTS
```

### Setup Inicial VPS

```bash
# 1. Login via SSH
ssh ubuntu@seu-vps.com

# 2. Atualizar sistema
sudo apt update && sudo apt upgrade -y

# 3. Instalar dependências base
sudo apt install -y curl wget git build-essential

# 4. Instalar NGINX
sudo apt install -y nginx

# 5. Instalar Certbot (Let's Encrypt)
sudo apt install -y certbot python3-certbot-nginx

# 6. Instalar golang-migrate
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# 7. Verificar instalações
nginx -v
```

---

## 🔄 NGINX

### nginx.conf

```nginx
upstream backend {
    server api:8080;
    keepalive 32;
}

server {
    listen 80;
    server_name api.seudominio.com;
    
    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.seudominio.com;
    
    # SSL certificates
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    
    # Security headers
    add_header Strict-Transport-Security \"max-age=31536000\" always;
    add_header X-Frame-Options \"SAMEORIGIN\" always;
    add_header X-Content-Type-Options \"nosniff\" always;
    
    # Compression
    gzip on;
    gzip_types text/plain text/css application/json;
    
    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=100r/s;
    limit_req zone=api burst=200 nodelay;
    
    # Proxy
    location / {
        proxy_pass http://backend;
        proxy_http_version 1.1;
        proxy_set_header Connection \"\";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # Health check
    location /health {
        access_log off;
        proxy_pass http://backend;
    }
}
```

### Let's Encrypt (Certbot)

```bash
# Gerar certificado
sudo certbot certonly --nginx -d api.seudominio.com

# Auto-renewal (já configurado)
sudo systemctl start certbot.timer
sudo systemctl enable certbot.timer

# Verificar renovação
sudo certbot renew --dry-run
```

---

## 🔄 CI/CD

### GitHub Actions Workflow

```yaml
# .github/workflows/deploy.yml
name: Build & Deploy

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.24.0
      
      - name: Run tests
        run: make test
      
      - name: Build binary
        run: make build
  
  deploy:
    needs: build
    runs-on: ubuntu-latest
    
    steps:
      - name: Deploy to VPS
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.VPS_HOST }}
          username: ${{ secrets.VPS_USER }}
          key: ${{ secrets.VPS_SSH_KEY }}
          script: |
            cd /opt/barber-api
            systemctl stop barber-api || true
            git pull
            make migrate-up
            make build
            systemctl start barber-api
```

---

## 📊 Monitoramento

### Prometheus Config

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'barber-api'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
```

### Alertas (AlertManager)

```yaml
# alertmanager.yml
route:
  receiver: 'backend-team'
  group_interval: 5m

receivers:
  - name: 'backend-team'
    email_configs:
      - to: 'ops@seudominio.com'
        from: 'alerts@seudominio.com'
        smarthost: 'smtp.gmail.com:587'
        auth_username: 'alerts@seudominio.com'
        auth_password: '${SMTP_PASSWORD}'
```

### Sentry Setup

```bash
# Backend
export SENTRY_DSN=https://xxxxx@sentry.io/xxxxx

# Frontend
# .env.local
NEXT_PUBLIC_SENTRY_DSN=https://xxxxx@sentry.io/xxxxx
```

---

## 💾 Backup & Recovery

### Backup Database (Neon)

```bash
# Automático via Neon (já configurado)

# Manual:
pg_dump -h neon.seudominio.com -U postgres -d barber_prod > backup.sql

# Restaurar:
psql -h neon.seudominio.com -U postgres -d barber_prod < backup.sql
```

### Backup Script (Cron)

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR=\"/backups\"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME=\"barber_prod\"

# Backup database
pg_dump $DATABASE_URL > $BACKUP_DIR/db_$DATE.sql

# Compress
gzip $BACKUP_DIR/db_$DATE.sql

# Upload para S3 (futuro)
# aws s3 cp $BACKUP_DIR/db_$DATE.sql.gz s3://backups/

# Deletar backups antigos (>30 dias)
find $BACKUP_DIR -name \"db_*.sql.gz\" -mtime +30 -delete

echo \"Backup completed at $DATE\"
```

Crontab:
```bash
# Rodar backup diariamente às 3h
0 3 * * * /opt/scripts/backup.sh
```

---

## 🔧 Troubleshooting

### API não responde

```bash
# 1. Verificar serviço systemd
sudo systemctl status barber-api

# 2. Verificar logs
journalctl -u barber-api -f

# 3. Verificar saúde
curl http://localhost:8080/health

# 4. Reiniciar
sudo systemctl restart barber-api
```

### Database connection error

```bash
# 1. Verificar DATABASE_URL
echo $DATABASE_URL

# 2. Testar conexão
psql $DATABASE_URL -c \"SELECT 1\"

# 3. Verificar migrations
migrate -path ./migrations -database $DATABASE_URL version

# 4. Rodar migrations
migrate -path ./migrations -database $DATABASE_URL up
```

### NGINX 502 Bad Gateway

```bash
# 1. Verificar status NGINX
sudo systemctl status nginx

# 2. Verificar config
sudo nginx -t

# 3. Verificar upstream
curl http://localhost:8080/health

# 4. Restart NGINX
sudo systemctl restart nginx
```

### High Memory Usage

```bash
# Identificar processo
ps aux | grep api

# Ver uso de memória
free -h
top -o %MEM

# Reiniciar serviço
sudo systemctl restart barber-api
```

---

## 📋 Checklist Deployment Produção

- [ ] SSL/TLS configurado e renovação automática
- [ ] NGINX rate limiting ativo
- [ ] Backups automáticos configurados
- [ ] Monitoramento e alertas ativos
- [ ] Logs centralizados (opcional: Sentry, Axiom)
- [ ] Health checks configurados
- [ ] CI/CD pipeline testado
- [ ] Rollback procedure documentado
- [ ] Postmortem template criado
- [ ] Documentação de runbooks completa

---

**Status:** ✅ Guia completo
