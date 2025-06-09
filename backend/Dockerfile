FROM golang:1.22-alpine

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Create data directory for SQLite
RUN mkdir -p /data

# Build the application
RUN go build -o main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/main"]
