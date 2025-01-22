# Gotham API

![Gotham Logo](https://res.cloudinary.com/dx70wyorg/image/upload/v1737406601/tentativa_do_gotham_utjkrf.png)

Gotham é uma API robusta desenvolvida em Go para gerenciamento de usuários, permissões e autenticação. O projeto utiliza PostgreSQL como banco de dados principal e Redis para cache de tokens, seguindo boas práticas de desenvolvimento e arquitetura modular.

## 🚀 Funcionalidades

- Autenticação JWT com refresh token
- Controle de acesso baseado em roles (RBAC)
- Cache de tokens com Redis
- Containerização com Docker
- CI/CD com GitHub Actions
- Migrations automáticas
- Seeds para dados iniciais
- Testes unitários e de integração

## 🛠️ Tecnologias

- **Go** 1.23+
- **PostgreSQL** 13
- **Redis** 6
- **Docker** e Docker Compose
- **JWT** para autenticação
- **GORM** como ORM
- **Gorilla Mux** para roteamento
- **GitHub Actions** para CI/CD
- **Testify** para testes

## 📦 Estrutura do Projeto

```
Gotham/
├── cmd/
│   └── gotham/
│       └── main.go
├── internal/
│   ├── database/
│   │   └── database.go
│   ├── handlers/
│   │   └── handlers.go
│   ├── middlewares/
│   │   └── middlewares.go
│   ├── models/
│   │   └── models.go
│   ├── routes/
│   │   └── routes.go
│   ├── seeds/
│   │   └── seeds.go
│   ├── settings/
│   │   └── config.go
│   └── utils/
│       └── utils.go
├── migrations/
│   └── migrations.go
├── pkg/
│   └── validator/
│       └── validator.go
├── test/
│   ├── config/
│   │   └── config.go
│   ├── helpers/
│   │   └── test_helper.go
│   ├── integration/
│   │   └── user_test.go
│   └── unit/
│       └── utils_test.go
├── .dockerignore
├── .env
├── .env-example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
└── README.md
```

## 🚀 Como Executar

1. Clone o repositório:
```bash
git clone https://github.com/jeffemart/Gotham.git
cd Gotham
```

2. Configure as variáveis de ambiente:
```bash
cp .env-example .env
# Edite o arquivo .env com suas configurações
```

3. Execute com Docker:
```bash
docker-compose up -d
```

4. Ou execute localmente:
```bash
go run cmd/gotham/main.go
```

## ⚡ Testes

Para executar os testes:
```bash
make test
```

## 🤝 Contribuição

1. Fork o projeto
2. Crie sua Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a Branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está licenciado sob a MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 👨‍💻 Autor

Jefferson Martins - [LinkedIn](https://www.linkedin.com/in/jefferson-martins-a6802b249/)
