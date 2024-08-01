# Use the official Golang image as a base
FROM golang:1.22

COPY app /app/app
COPY common /app/common
COPY config /app/config
COPY static /app/static
COPY errlog /app/errlog
COPY main.go /app
COPY go.mod /app
COPY go.sum /app
COPY Makefile /app

COPY image/gpjc_scripts /app/gpjc_scripts
COPY ca.crt /app
COPY ca.key /app

COPY wrapper-api /app/wrapper-api

# Copy the entrypoint script into the container
COPY image/entrypoint.sh /usr/local/bin/entrypoint.sh

# Make the entrypoint script executable
RUN chmod +x /usr/local/bin/entrypoint.sh

# Set the working directory inside the container
WORKDIR /app

# Accept build argument
ARG DNS_NAME
ARG MPC_IMPLEMENTATION

# Set the environment variable
ENV DNS_NAME=${DNS_NAME}
ENV MPC_IMPLEMENTATION=${MPC_IMPLEMENTATION}

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Expose the port that the app listens on
EXPOSE 4000

# Default command, will be run if not overridden in 'docker run'
CMD go run .
