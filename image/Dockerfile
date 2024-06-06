# Use the official Golang image as a base
FROM golang:1.19

COPY app /app/app
COPY common /app/common
COPY config /app/config
COPY static /app/static
COPY main.go /app
COPY go.mod /app
COPY go.sum /app
COPY Makefile /app

COPY image/gpjc_scripts /app/gpjc_scripts
COPY ca.crt /app
COPY ca.key /app

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
RUN chmod +x /app/gpjc_scripts/certs_script.sh

# Run the scripts
RUN /app/gpjc_scripts/certs_script.sh

RUN mv /app/ca.crt /app/private-join-and-compute && \
    mv /app/ca.key /app/private-join-and-compute && \
    mv /app/client.crt /app/private-join-and-compute && \
    mv /app/client.csr /app/private-join-and-compute && \
    mv /app/client.key /app/private-join-and-compute && \
    mv /app/server.crt /app/private-join-and-compute && \
    mv /app/server.csr /app/private-join-and-compute && \
    mv /app/server.key /app/private-join-and-compute

# Expose the port that the app listens on
EXPOSE 4000

# Run the application
CMD /app/gpjc-api ${DNS_NAME} Ethernal123 & go run .