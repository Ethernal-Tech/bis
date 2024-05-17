#!/bin/bash

# Function to check if BIS database exists
check_db_exists() {
  /opt/mssql-tools/bin/sqlcmd -S sql-database -U SA -P Ethernal123 -Q "IF DB_ID('BIS') IS NOT NULL SELECT 1 AS db_exists ELSE SELECT 0 AS db_exists" -h -1 -W | grep -q 1
}

# Wait until MSSQL is ready to accept connections
until nc -z -v -w30 sql-database 1433; do
  echo "Waiting for MSSQL to be ready..."
  sleep 10
done

# Wait until the BIS database is created
until check_db_exists; do
  echo "Waiting for BIS database to be created..."
  sleep 10
done

# MSSQL is ready and BIS database is created, start the application
exec "$@"
