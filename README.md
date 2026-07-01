# Getting Started with Golang App

This project utilizes the [Gin-Gonic](https://github.com/gin-gonic/gin) framework in Golang

# Project Dependencies

- **GO Version**: 1.21.5

## External Libraries

1. **gin-gonic/gin v1.9.1**
   - Framework for building HTTP web services in Go.
   
2. **gin-contrib/cors v1.5.0**
   - Middleware for enabling Cross-Origin Resource Sharing (CORS) in Gin.
   
3. **golang-jwt/jwt/v5 v5.2.0**
   - JSON Web Token (JWT) implementation for Go.
   
4. **minio/minio-go/v7 v7.0.66**
   - MinIO Go SDK for interacting with the MinIO object storage server.
   
5. **gorm.io/gorm v1.25.5**
   - Object Relational Mapping (ORM) library for Golang.
   
6. **githubnemo/CompileDaemon v1.4.0**
   - Compile Daemon for Golang, automatically recompiles and restarts your Go application on code changes.

## Installation

Use the go get to download and install all the packages used by your project

```bash
go get -u all
```

Install Compile Daemon for run the project

```bash
go install github.com/githubnemo/CompileDaemon
```

## Running the Project

### Using 'go run'

```bash
go run main.go
```

### Using Compile Daemon

```bash
compiledaemon --command="./book-recipe-be-go"
```

### Using Docker

#### Build Image

```bash
docker build --tag book-recipe-be-go .
```

#### Build Network

```bash
docker network create book-recipe-be-go-network
```

#### Build Container and Run

```bash
docker run -d -p 8080:8080 --name book-recipe-be-go --network book-recipe-be-go-network book-recipe-be-go
```

# URL Health Check

```bash
http://localhost:8080/api/health
```

# Learn More

You can learn more in the [Golang Documentation](https://go.dev/doc/)

To learn Gin-Gonic Framework, check out the [Gin-Gonic Documentation](https://gin-gonic.com/docs/)

To learn GORM, check out the [GORM Documentation](https://gorm.io/docs/)
