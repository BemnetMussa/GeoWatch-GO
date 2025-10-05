# Stage 1: Build the Go binary
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application. CGO_ENABLED=0 is important for a static binary.
# AFTER - The correct path, relative to the /app workdir
RUN CGO_ENABLED=0 GOOS=linux go build -o /geowatch-backend ./cmd/main.go

# Stage 2: Create the final, lightweight image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the 'builder' stage
COPY --from=builder /geowatch-backend .

# Expose port 8000 (or whatever port your Gin server uses)
EXPOSE 8000

# Command to run the executable
CMD ["./geowatch-backend"]