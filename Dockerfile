FROM golang:1.23

# Definir diretório de trabalho no contêiner
WORKDIR /app

# Copiar todo o código fonte
COPY . .

# Instalar dependências e compilar
RUN go mod download && \
    go mod tidy && \
    go build -o /usr/local/bin/app ./cmd/gotham/main.go && \
    chmod +x /usr/local/bin/app

# Criar diretórios necessários
RUN mkdir -p /app/docs && \
    mkdir -p /app/api/swagger

# Copiar arquivos da documentação
COPY docs/openapi.json /app/docs/
COPY api/swagger/* /app/api/swagger/

# Expor a porta 8000
EXPOSE 8000

# Rodar a aplicação
CMD ["/usr/local/bin/app"]
