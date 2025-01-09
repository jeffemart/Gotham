# Gotham

Gotham é uma aplicação construída em Go, utilizando PostgreSQL como banco de dados e Redis como cache. Este projeto segue uma arquitetura modular e utiliza Docker para gerenciar os serviços e facilitar a implantação.

## Estrutura do Projeto

```
Gotham/
├── app/
│   ├── handlers/
│   ├── middlewares/
│   ├── models/
│   ├── routes/
│   ├── utils/
│   └── main.go
├── docker-compose.yml
├── .env
└── README.md
```

- **handlers/**: Contém as funções para lidar com as rotas da API.
- **middlewares/**: Middleware para validação de autenticação e autorização.
- **models/**: Definições de modelos e interações com o banco de dados.
- **routes/**: Configuração de rotas e agrupamentos.
- **utils/**: Funções utilitárias, como parsing de JWT.
- **main.go**: Ponto de entrada da aplicação.

## Pré-requisitos

Certifique-se de ter os seguintes itens instalados:

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- Go 1.20+

## Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com o seguinte conteúdo:

```
# Variáveis para a aplicação Go
APP_PORT=8000

# Variáveis para o Postgres
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=gotham_db
POSTGRES_PORT=5432

# Variáveis para o Redis
REDIS_HOST=redis
REDIS_PORT=6379
```

### Subindo os Serviços

1. Clone o repositório:

   ```bash
   git clone https://github.com/jeffemart/Gotham.git
   cd Gotham
   ```

2. Suba os serviços utilizando Docker Compose:

   ```bash
   docker-compose --env-file .env up --build
   ```

3. Acesse a aplicação em `http://localhost:8000`.

## Endpoints da API

### Rotas Públicas

- `GET /users`: Retorna a lista de usuários.
- `GET /users/{id}`: Retorna um usuário específico pelo ID.

### Rotas Protegidas (Admin)

- `POST /admin/users`: Cria um novo usuário.
- `PUT /admin/users/{id}`: Atualiza um usuário existente.
- `DELETE /admin/users/{id}`: Exclui um usuário.

### Rotas Protegidas (Admin ou Agente)

- `GET /protected/tasks`: Retorna a lista de tarefas.

### Autenticação

- `POST /login`: Gera um token JWT para autenticação.
- `POST /refresh`: Gera um novo token JWT a partir de um token válido.

## Testando a API

Utilize ferramentas como [Postman](https://www.postman.com/) ou [cURL](https://curl.se/) para testar os endpoints.

### Exemplo de Requisição

#### Login

```bash
curl -X POST http://localhost:8000/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin"}'
```

#### Criar Usuário (Admin)

```bash
curl -X POST http://localhost:8000/admin/users \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "johndoe@example.com"}'
```

## Tecnologias Utilizadas

- **Linguagem**: Go
- **Banco de Dados**: PostgreSQL
- **Cache**: Redis
- **Containerização**: Docker

## Licença

Este projeto está sob a licença [MIT](https://opensource.org/licenses/MIT).

