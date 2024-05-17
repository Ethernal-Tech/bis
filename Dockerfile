# Use the official Golang image as a base
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project into the container
COPY . .

# Install dependencies including mssql-tools and netcat
RUN apt-get update && \
    apt-get install -y make curl gnupg2 apt-transport-https software-properties-common && \
    curl https://packages.microsoft.com/keys/microsoft.asc | apt-key add - && \
    curl https://packages.microsoft.com/config/ubuntu/20.04/prod.list | tee /etc/apt/sources.list.d/msprod.list && \
    apt-get update && \
    ACCEPT_EULA=Y apt-get install -y mssql-tools unixodbc-dev netcat && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Add SQLCMD to PATH
ENV PATH="$PATH:/opt/mssql-tools/bin"

# # Fetch the necessary releases using make
# RUN make fetch-releases

# Copy the wait-for-db.sh script
COPY wait-for-db.sh /app/wait-for-db.sh

# Make the script executable
RUN chmod +x /app/wait-for-db.sh

# Expose the port that the app listens on
EXPOSE 4000

# Run the application using the wait-for-db.sh script
CMD ["/app/wait-for-db.sh", "go", "run", "."]
