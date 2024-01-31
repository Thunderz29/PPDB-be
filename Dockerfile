# Stage 1: Build the Go application
FROM golang:1.21 AS builder

WORKDIR /usr/src/app

# Copy the entire project to the working directory
COPY . .
RUN go mod download
RUN go mod verify
# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /book-recipe-be-go .

# Stage 2: Create a minimal image to run the application
FROM alpine:3.16

WORKDIR /app

# Copy the .env file
COPY .env .

RUN apk update

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy the executable.
COPY --from=builder /book-recipe-be-go /book-recipe-be-go

ENTRYPOINT ["/book-recipe-be-go"]
