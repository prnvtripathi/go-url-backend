# Step 1: Use the official Go image to build the app
FROM golang:1.23.0-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies (this will be cached if go.mod and go.sum haven't changed)
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Step 2: Run the app with a minimal Docker image
FROM alpine:latest

# Set working directory in container
WORKDIR /root/

# Copy the pre-built binary file from the builder stage
COPY --from=builder /app/main .

# Copy any necessary environment files
COPY .env .

# Expose port (change if your server listens on another port)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
