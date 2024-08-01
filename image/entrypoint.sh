#!/bin/bash

# Check if a specific condition is met (e.g., if VAR is not equal to SL)
if [ "$MPC_IMPLEMENTATION" != "SL" ]; then
    echo "Using google/private-join-and-compute MPC implementation"

    # Fetch the necessary releases using make
    make fetch-releases-multiple-machines

    # Set the DNS_NAME environment variable if it was passed as an argument
    export DNS_NAME=${DNS_NAME}

    # Make the script executable
    chmod +x /app/gpjc_scripts/certs_script.sh

    # Run the certs script
    /app/gpjc_scripts/certs_script.sh

    # Move the generated certificates to the specified directory
    mv /app/ca.crt /app/private-join-and-compute
    mv /app/ca.key /app/private-join-and-compute
    mv /app/client.crt /app/private-join-and-compute
    mv /app/client.csr /app/private-join-and-compute
    mv /app/client.key /app/private-join-and-compute
    mv /app/server.crt /app/private-join-and-compute
    mv /app/server.csr /app/private-join-and-compute
    mv /app/server.key /app/private-join-and-compute

    # Start the gpjc-api with arguments
    /app/gpjc-api ${DNS_NAME} Ethernal123 &
else
    echo "Using SL MPC implementation"
    
    /app/wrapper-api &
fi

# Optionally, you can start the main process of your container here, e.g., a web server or application
exec "$@"
