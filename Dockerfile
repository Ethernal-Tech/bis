# Use the official Golang image as a base
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project into the container
COPY . .

# Add SQLCMD to PATH
ENV PATH="$PATH:/opt/mssql-tools/bin"

# # Fetch the necessary releases using make
# RUN make fetch-releases

# Expose the port that the app listens on
EXPOSE 4000

# Run the application using the wait-for-db.sh script
CMD ["go", "run", "."]
