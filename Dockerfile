FROM golang:1.23

# Definir diretório de trabalho no contêiner
WORKDIR /app

# Copiar arquivos de dependência primeiro
COPY go.mod go.sum ./

# Baixar dependências e atualizar go.sum
RUN go mod download && \
    go mod tidy

# Copiar o resto do código fonte
COPY . .

# Compilar a aplicação
RUN go build -o /usr/local/bin/app ./cmd/gotham/main.go && \
    chmod +x /usr/local/bin/app

# Expor a porta 8000
EXPOSE 8000

# Rodar a aplicação
CMD ["/usr/local/bin/app"]
