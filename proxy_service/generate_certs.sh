#!/usr/bin/env bash
set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <BACKEND_IP>"
  exit 1
fi

BACKEND_IP="$1"
CERT_DIR="./certs"

mkdir -p "$CERT_DIR"
cd "$CERT_DIR"

echo "=== Using Backend IP: $BACKEND_IP ==="

echo "=== 1. Generate Root CA ==="

openssl genrsa -out ca.key 4096

openssl req -x509 -new -nodes \
  -key ca.key \
  -sha256 -days 3650 \
  -out ca.pem \
  -subj "/C=IR/O=Internal-CA/CN=Internal Root CA"

echo "=== 2. Generate Backend Server Certificate ==="

openssl genrsa -out backend-server.key 4096

openssl req -new \
  -key backend-server.key \
  -out backend-server.csr \
  -subj "/C=IR/O=Backend/CN=backend.local"

cat > backend-ext.cnf <<EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
IP.1 = ${BACKEND_IP}
DNS.1 = backend.local
EOF

openssl x509 -req \
  -in backend-server.csr \
  -CA ca.pem \
  -CAkey ca.key \
  -CAcreateserial \
  -out backend-server.pem \
  -days 825 \
  -sha256 \
  -extfile backend-ext.cnf

echo "=== 3. Generate Traefik Client Certificate ==="

openssl genrsa -out traefik-client.key 4096

openssl req -new \
  -key traefik-client.key \
  -out traefik-client.csr \
  -subj "/C=IR/O=Traefik/CN=traefik-client"

cat > traefik-client-ext.cnf <<EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth
EOF

openssl x509 -req \
  -in traefik-client.csr \
  -CA ca.pem \
  -CAkey ca.key \
  -CAcreateserial \
  -out traefik-client.pem \
  -days 825 \
  -sha256 \
  -extfile traefik-client-ext.cnf

echo "=== 4. Cleanup ==="
rm -f *.csr *.cnf *.srl

echo "=== DONE ==="
echo
echo "Generated files:"
ls -1
