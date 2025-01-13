# Etapa 1: Build da aplicação Go
FROM golang:1.23 AS build

# Definir diretório de trabalho no contêiner
WORKDIR /app

# Copiar o código da aplicação para o contêiner
COPY . .

# Instalar dependências e compilar o código Go
RUN go mod tidy
RUN go build -o main .

# Etapa 2: Execução da aplicação
FROM golang:1.23

# Definir diretório de trabalho no contêiner
WORKDIR /app

# Copiar o binário da aplicação da etapa de build
COPY --from=build /app/main .

# Expor a porta 8000
EXPOSE 8000
EXPOSE 8080

# Rodar a aplicação Go
CMD ["./main"]
