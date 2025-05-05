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

# Check if the database already exists by querying the pg_database catalog.
DB_EXISTS=$(psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -tAc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}';")

if [ "$DB_EXISTS" = "1" ]; then
    echo "Database '${DB_NAME}' already exists."
else
    echo "Database '${DB_NAME}' does not exist. Creating it..."
    psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -c "CREATE DATABASE ${DB_NAME};"

    if [ $? -eq 0 ]; then
        echo "Database '${DB_NAME}' created successfully."
    else
        echo "Error creating database '${DB_NAME}'."
    fi
fi

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
INSERT INTO users (username, password, role_id)
VALUES ('${ADMIN_USERNAME}', '${ADMIN_PASSWORD_HASH}', 1)
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
