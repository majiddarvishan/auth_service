# Auth service

## Get an SSL Certificate

1. Use a Self-Signed Certificate (for testing)
2. Obtain a Certificate from Let’s Encrypt (for production)

### Generate a Self-Signed Certificate

```bash
openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out cert.pem -days 365
```

### Using Let’s Encrypt

1. Install Certbot

```bash
sudo apt install certbot
```

2. Get an SSL certificate

```bash
certbot certonly --standalone -d yourdomain.com
```


Use mkcert (better!)

```bash
# Install mkcert (Linux example)
sudo apt install libnss3-tools
sudo apt install mkcert
mkcert -install
mkcert localhost 127.0.0.1 ::1
```

## Generate the Swagger Documentation

```sh
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files

swag init -g main.go  -o docs
```

- to view output

```bash
https://<URL>/swagger/index.html
```

### swaggertodoc

https://mvnrepository.com/artifact/io.github.swagger2markup/swagger2markup-cli
wget https://repo1.maven.org/maven2/io/github/swagger2markup/swagger2markup-cli/1.3.3/swagger2markup-cli-1.3.3.jar
java -jar /path/to/swagger2markup-cli-1.3.3.jar convert -i <your_swagger_file.json_or_yaml> -d <output_directory>

swagger2markup convert -i ./docs/swagger.json -f output.adoc
pandoc output.adoc -o output.docx


pip install python-docx
python swagger2docx.py

## Running with Docker Compose

### Prerequisites

- Docker and Docker Compose installed
- TLS certificates in `deployments/tls/` directory:
  - `localhost.pem` (certificate)
  - `localhost-key.pem` (private key)

### Setup and Run

1. Generate TLS certificates (if not present):

```bash
cd deployments
mkcert localhost 127.0.0.1 ::1
mv localhost.pem tls/
mv localhost-key.pem tls/
cd ..
```

2. Build the Docker image:

```bash
docker-compose -f deployments/docker-compose.yml build
```

3. Start services:

```bash
docker-compose -f deployments/docker-compose.yml up -d
```

4. Verify services are running:

```bash
docker-compose -f deployments/docker-compose.yml ps
```

### Environment Variables (in docker-compose.yml)

| Variable | Value | Notes |
|----------|-------|-------|
| `GIN_MODE` | `release` | Production mode (no debug logs) |
| `BASE_API` | `/v1/api` | API endpoint prefix |
| `TLS_PATH` | `/etc/tls` | Path to TLS certs inside container |
| `SECRET_KEY` | 64-char hex | JWT signing key (generate: `openssl rand -hex 32`) |
| `TOKEN_EXPIRATION_PERIOD` | `24h` | Access token lifetime |
| `DB_HOST` | `postgres` | PostgreSQL host (use service name) |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | PostgreSQL user |
| `DB_PASSWORD` | `postgres` | PostgreSQL password |
| `DB_NAME` | `proxy_db` | Database name |
| `ACCOUNTING_ENDPOINT` | `http://...` | External accounting service (optional) |

### Service Health

Both services include health checks:
- **auth_proxy**: Listens on HTTPS 8443
- **postgres**: PostgreSQL healthcheck with pg_isready

### View Logs

```bash
# All services
docker-compose -f deployments/docker-compose.yml logs -f

# Specific service
docker-compose -f deployments/docker-compose.yml logs -f auth_proxy
docker-compose -f deployments/docker-compose.yml logs -f postgres
```

### Access API

```bash
# Health check
curl -k https://localhost:8443/health

# Swagger UI
# Open browser to https://localhost:8443/swagger/index.html
```

### Cleanup

```bash
# Stop services
docker-compose -f deployments/docker-compose.yml down

# Remove volumes (wipe database)
docker-compose -f deployments/docker-compose.yml down -v
```

## Install PostgreSQL

```bash
apt-get install postgresql
```

## Install mkpasswd

```bash
apt-get install whois
```

## Generate JWT_SECRET

```bash
openssl rand -hex 32
```
