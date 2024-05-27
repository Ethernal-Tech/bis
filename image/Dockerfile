# Use the official Golang image as a base
FROM golang:1.19

COPY common /app/common
COPY config /app/config
COPY controller /app/controller
COPY core /app/core
COPY server /app/server
COPY static /app/static
COPY main.go /app
COPY go.mod /app
COPY go.sum /app
COPY Makefile /app

COPY image/gpjc_scripts /app/gpjc_scripts

# Set the working directory inside the container
WORKDIR /app

# Fetch the necessary releases using make
RUN make fetch-releases-multiple-machines

# Install ODBC library for GPJC API
RUN apt-get update && apt-get install -y unixodbc unixodbc-dev

# Accept build argument
ARG DNS_NAME

# Set the environment variable
ENV DNS_NAME=${DNS_NAME}

# Make the scripts executable
RUN chmod +x /app/gpjc_scripts/ca_script.sh
RUN chmod +x /app/gpjc_scripts/certs_script.sh

# Run the scripts
RUN /app/gpjc_scripts/ca_script.sh
RUN /app/gpjc_scripts/certs_script.sh

# Expose the port that the app listens on
EXPOSE 4000

# Run the application
CMD /app/gpjc-api app1 Ethernal123 & go run .