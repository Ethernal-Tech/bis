# Use the official Golang image to build the application
FROM golang:1.22.3-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o /go/bin/p2pmock

# Use a smaller image for the runtime
FROM alpine:latest

# Copy the compiled Go binary from the builder stage
COPY --from=builder /go/bin/p2pmock /p2pmock

# Expose the port the app runs on
EXPOSE 5000

# Command to run the executable
CMD ["/p2pmock"]
