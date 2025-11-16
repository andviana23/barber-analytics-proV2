# 🚀 Guia DevOps

**Versão:** 1.0  
**Data:** 14/11/2025  
**Status:** Guia Prático

---

## 📋 Índice

1. [Infraestrutura](#infraestrutura)
2. [Docker & Compose](#docker--compose)
3. [NGINX](#nginx)
4. [CI/CD](#cicd)
5. [Monitoramento](#monitoramento)
6. [Backup & Recovery](#backup--recovery)
7. [Troubleshooting](#troubleshooting)

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

# 4. Instalar Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 5. Adicionar user ao docker group
sudo usermod -aG docker ubuntu

# 6. Instalar Docker Compose
sudo curl -L \"https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)\" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 7. Instalar NGINX
sudo apt install -y nginx

# 8. Instalar Certbot (Let's Encrypt)
sudo apt install -y certbot python3-certbot-nginx

# 9. Instalar golang-migrate
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# 10. Verificar instalações
docker --version
docker-compose --version
nginx -v
```

---

## 🐳 Docker & Compose

### Dockerfile Backend

```dockerfile
# Dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/api .

EXPOSE 8080

CMD [\"./api\"]
```

### docker-compose.yml (Dev)

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - \"8080:8080\"
    environment:
      DATABASE_URL: postgres://postgres:password@db:5432/barber_dev
      JWT_SECRET: dev-secret-key
      APP_ENV: development
    depends_on:
      - db
    volumes:
      - .:/app
    command: go run ./cmd/api/main.go

  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: barber_dev
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - \"5432:5432\"

volumes:
  postgres_data:
```

### Docker Compose (Produção)

```yaml
version: '3.8'

services:
  api:
    image: seu-registry/barber-api:latest
    restart: always
    environment:
      DATABASE_URL: ${DATABASE_URL}
      JWT_SECRET: ${JWT_SECRET}
      APP_ENV: production
    ports:
      - \"8080:8080\"
    healthcheck:
      test: [\"CMD\", \"curl\", \"-f\", \"http://localhost:8080/health\"]
      interval: 30s
      timeout: 10s
      retries: 3

  nginx:
    image: nginx:alpine
    restart: always
    ports:
      - \"80:80\"
      - \"443:443\"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - api
```

### Comandos Úteis

```bash
# Build imagem
docker build -t barber-api:v2.0 .

# Rodar container
docker run -d --name barber-api -p 8080:8080 barber-api:v2.0

# Ver logs
docker logs -f barber-api

# Parar container
docker stop barber-api

# Docker Compose
docker-compose up -d
docker-compose logs -f api
docker-compose down
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
          go-version: 1.22
      
      - name: Run tests
        run: make test
      
      - name: Build Docker image
        run: docker build -t seu-registry/barber-api:${{ github.sha }} .
      
      - name: Push to registry
        run: |
          echo \"${{ secrets.REGISTRY_PASSWORD }}\" | docker login -u ${{ secrets.REGISTRY_USER }} --password-stdin
          docker push seu-registry/barber-api:${{ github.sha }}
  
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
            cd /app/barber-analytics
            docker-compose pull
            docker-compose up -d
            docker-compose exec -T api migrate -path ./migrations -database $DATABASE_URL up
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
# 1. Verificar container
docker ps | grep barber-api

# 2. Verificar logs
docker logs barber-api

# 3. Verificar saúde
curl http://localhost:8080/health

# 4. Reiniciar
docker restart barber-api
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
docker stats

# Limpar cache
docker system prune

# Restart container
docker-compose restart api
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
