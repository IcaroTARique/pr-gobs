# Usar uma imagem base do Golang para construir a aplicação
FROM golang:1.21-alpine as builder

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos go.mod e go.sum e baixa as dependências
COPY go.mod go.sum ./
COPY . .
COPY ./cmd/server/.env .

ENV GOPROXY="https://goproxy.io"
RUN go mod tidy


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server/main.go


FROM alpine:latest

COPY --from=builder /app/server .
COPY ./cmd/server/.env .

EXPOSE 8081

CMD ["./server"]
