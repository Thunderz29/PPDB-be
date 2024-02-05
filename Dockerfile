# Stage 1: Build the Go application
FROM harbor.cloudias79.com/devops-tools/golang:1.21-alpine AS builder

WORKDIR /app

# Copy the entire project to the working directory
COPY . .
RUN go mod download
RUN go mod verify
# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/book-recipe-be-go .

# Stage 2: Create a minimal image to run the application
FROM harbor.cloudias79.com/devops-tools/alpine:3.14.4
RUN adduser -D -g '' appuser
USER appuser

WORKDIR /app

# Copy the .env file
COPY --chown=appsuser:appuser .env /app/.env
COPY --chown=appsuser:appuser app.log /app/app.log
# Copy the executable.
COPY --chown=appsuser:appuser --from=builder /app/book-recipe-be-go /app/book-recipe-be-go
# RUN chmod +x /app/book-recipe-be-go

CMD ["/app/book-recipe-be-go"]
