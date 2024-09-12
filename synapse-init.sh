#!/bin/sh
set -e

# Check if config file exists
if [ ! -f /data/homeserver.yaml ]; then
    echo "Generating Synapse config..."
    python -m synapse.app.homeserver \
        --server-name $SYNAPSE_SERVER_NAME \
        --config-path /data/homeserver.yaml \
        --generate-config \
        --report-stats=no
    
    # Modify the generated config to use PostgreSQL
    sed -i 's/database: "sqlite3"/database: "psycopg2"/' /data/homeserver.yaml
    echo "
database:
    name: psycopg2
    args:
        user: $POSTGRES_USER
        password: $POSTGRES_PASSWORD
        database: $POSTGRES_DB
        host: postgres
        cp_min: 5
        cp_max: 10
" >> /data/homeserver.yaml
fi

# Start Synapse
echo "Starting Synapse..."
python -m synapse.app.homeserver --config-path /data/homeserver.yaml
