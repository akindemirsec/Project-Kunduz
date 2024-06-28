FROM golang:1.19-alpine

WORKDIR /app

# Copy go.mod and go.sum files first
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

RUN apk add --no-cache python3 py3-pip
RUN pip3 install requests psycopg2-binary

# Build the application
RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]
