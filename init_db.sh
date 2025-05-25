#!/bin/bash
DB_EXISTS=$(psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname='${APP_NAME}db'" | grep -q 1 && echo "yes" || echo "no")

if [ "$DB_EXISTS" == "no" ]; then
    psql -U postgres -c "CREATE DATABASE ${APP_NAME}db;"
fi
