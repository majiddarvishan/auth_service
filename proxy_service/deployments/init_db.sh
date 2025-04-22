#!/bin/bash
set -e

# apt-get install apache2-utils
# apt-get install postgresql-client

# ------------------------------
# Configuration Section
# ------------------------------

# Database configuration
DB_HOST="172.23.10.73"
DB_PORT="5432"
DB_USER="postgres"
DB_PASSWORD="postgres"
DB_NAME="mydb"

# Admin user configuration
ADMIN_USERNAME="admin"
ADMIN_PASSWORD="verysecurepassword"  # Use a secure method to provide/manage this

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
# SQL: Insert Admin Role
# ------------------------------
# This assumes your roles table has a unique constraint on `name`
ROLE_SQL="
INSERT INTO roles (name, description)
VALUES ('admin', 'admin role with all permissions')
ON CONFLICT (name) DO NOTHING;
"

psql \
  -h "${DB_HOST}" \
  -p "${DB_PORT}" \
  -U "${DB_USER}" \
  -d "${DB_NAME}" \
  -c "${ROLE_SQL}"

echo "Admin role ensured in the database."

# ------------------------------
# SQL Query to Ensure Default Admin User
# ------------------------------
SQL_QUERY="
INSERT INTO users (username, password, role)
VALUES ('${ADMIN_USERNAME}', '${ADMIN_PASSWORD_HASH}', 'admin')
ON CONFLICT (username) DO NOTHING;
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

echo "Default admin user ensured in the database."
