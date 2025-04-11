# Use official Go image based on Alpine
FROM golang:1.24.2-alpine3.21

# Install system tools: make, git, bash, curl, unzip
RUN apk add --no-cache make git bash curl unzip

# Install golang-migrate CLI
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz \
    -o migrate.tar.gz \
    && tar -xzf migrate.tar.gz \
    && mv migrate /usr/local/bin/migrate \
    && rm -f migrate.tar.gz

# Set working directory inside container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy entire project directory (source files, Makefile, .env, migrations, etc.)
COPY . .

# Build Go binary
RUN go build -o main .

# Command to run the binary
CMD ["./main"]
