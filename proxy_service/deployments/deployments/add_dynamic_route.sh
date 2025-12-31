#!/bin/bash
set -e

# apt-get install apache2-utils
# apt-get install postgresql-client

# ------------------------------
# Configuration Section
# ------------------------------

# Database configuration
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="postgres"
DB_PASSWORD="postgres"
DB_NAME="mydb"

# Admin user configuration
ADMIN_USERNAME="admin"
ADMIN_PASSWORD="admin"  # Use a secure method to provide/manage this

# ------------------------------
# Generate Bcrypt Hash for the Admin Password
# ------------------------------
# First, try mkpasswd if available, then htpasswd, otherwise fall back
if command -v mkpasswd >/dev/null 2>&1; then
    ADMIN_PASSWORD_HASH=$(mkpasswd -m bcrypt "${ADMIN_PASSWORD}")
elif command -v htpasswd >/dev/null 2>&1; then
    # htpasswd returns: "username:hash", so extract the hash portion.
    ADMIN_PASSWORD_HASH=$(htpasswd -nbB admin "${ADMIN_PASSWORD}" | cut -d ':' -f2)
else
    ADMIN_PASSWORD_HASH="${ADMIN_PASSWORD}"
    echo "Warning: Neither mkpasswd nor htpasswd were found. The admin password will be stored in plaintext. This is not safe for production."
fi

# Export the password so that psql can use it
export PGPASSWORD="${DB_PASSWORD}"

# ------------------------------
# SQL Query to Route SMS
# ------------------------------

SQL_QUERY="
INSERT INTO custom_endpoints (created_at, updated_at, path, method, endpoints, need_accounting, enabled)
VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '/sms/*path', 'ANY', '{http://192.168.1.225:8080/api/v1}', false, true)
ON CONFLICT (path) DO NOTHING;
"

# ------------------------------
# Execute the SQL Query Using psql
# ------------------------------
psql \
  -h "${DB_HOST}" \
  -p "${DB_PORT}" \
  -U "${DB_USER}" \
  -d "${DB_NAME}" \
  -c "${SQL_QUERY}"

# ------------------------------
# SQL Query to Route prefixes
# ------------------------------

SQL_QUERY="
INSERT INTO custom_endpoints (created_at, updated_at, path, method, endpoints, need_accounting, enabled)
VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '/prefix/*path', 'ANY', '{http://192.168.1.225:8080/api/v1}', false, true)
ON CONFLICT (path) DO NOTHING;
"

# ------------------------------
# Execute the SQL Query Using psql
# ------------------------------
psql \
  -h "${DB_HOST}" \
  -p "${DB_PORT}" \
  -U "${DB_USER}" \
  -d "${DB_NAME}" \
  -c "${SQL_QUERY}"

echo "Dynamic routes are added in the database."
