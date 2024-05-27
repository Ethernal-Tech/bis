# Use the official Golang image as a base
FROM golang:1.19

COPY config /app/config
COPY controller /app/controller
COPY core /app/core
COPY server /app/server
COPY static /app/static
COPY main.go /app
COPY go.mod /app
COPY go.sum /app


# Set the working directory inside the container
WORKDIR /app

# # Fetch the necessary releases using make
# RUN make fetch-releases

# Expose the port that the app listens on
EXPOSE 4000

# Run the application using the wait-for-db.sh script
CMD ["go", "run", "."]