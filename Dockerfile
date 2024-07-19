FROM golang:1.22.5 AS build-stage
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmp/main.go

# for production test the use from scratch

EXPOSE 8080

ENTRYPOINT ["/api"]
