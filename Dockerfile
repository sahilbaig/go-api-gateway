# ---- Build stage ----
FROM golang:1.24-alpine AS build

WORKDIR /app

# Copy go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the Go binary
RUN go build -o app .

# ---- Run stage ----
FROM alpine:latest

WORKDIR /app

# Copy the binary from build stage
COPY --from=build /app/app .

# Expose the port
EXPOSE 8080

CMD ["./app"]
