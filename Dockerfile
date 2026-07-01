# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Copy the entire project to the working directory
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/book-recipe-be-go .

# Stage 2: Create a minimal image to run the application
FROM alpine:3.18
RUN adduser -D -g '' appuser
USER appuser

WORKDIR /app

RUN touch app.log

# Copy the .env file
COPY --chown=appuser:appuser --from=builder /app/book-recipe-be-go /app/book-recipe-be-go
# RUN chmod +x /app/book-recipe-be-go

EXPOSE 8080

CMD ["/app/book-recipe-be-go"]