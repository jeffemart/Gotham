# Gotham API

Gotham é uma API robusta desenvolvida em Go para gerenciamento de usuários, permissões e autenticação. O projeto utiliza PostgreSQL como banco de dados principal e Redis para cache de tokens, seguindo boas práticas de desenvolvimento e arquitetura modular.

## 🚀 Funcionalidades

- Autenticação JWT com refresh token
- Controle de acesso baseado em roles (RBAC)
- Cache de tokens com Redis
- Documentação automática com Swagger
- Containerização com Docker
- CI/CD com GitHub Actions

## 🛠️ Tecnologias

- **Go** 1.23+
- **PostgreSQL** 13
- **Redis** 6
- **Docker** e Docker Compose
- **JWT** para autenticação
- **GORM** como ORM
- **Swagger** para documentação
- **GitHub Actions** para CI/CD

## 📦 Estrutura do Projeto

```
Gotham/
├── database/       # Configuração e conexão com bancos de dados
├── docs/          # Documentação Swagger
├── handlers/      # Handlers HTTP
├── middlewares/   # Middlewares de autenticação e autorização
├── migrations/    # Migrações do banco de dados
├── models/        # Modelos de dados
├── routes/        # Configuração de rotas
├── settings/      # Configurações da aplicação
├── utils/         # Funções utilitárias
└── main.go        # Ponto de entrada da aplicação
```

## 🚦 Endpoints da API

### Rotas Públicas
- `POST /login` - Autenticação de usuário
- `GET /users` - Lista todos os usuários
- `GET /users/{id}` - Obtém um usuário específico
- `POST /users` - Cria um novo usuário

### Rotas Protegidas (Admin)
- `PUT /admin/users/{id}` - Atualiza um usuário
- `DELETE /admin/users/{id}` - Remove um usuário

### Rotas Protegidas (Admin ou Agente)
- `GET /protected/tasks` - Lista tarefas

## 🛠️ Instalação

### Pré-requisitos
- Docker e Docker Compose
- Go 1.23+
- Make (opcional)

### Configuração

1. Clone o repositório:
```bash
git clone https://github.com/jeffemart/Gotham.git
cd Gotham
```

2. Copie o arquivo de exemplo de ambiente:
```bash
cp .env-example .env
```

3. Configure as variáveis de ambiente no arquivo `.env`

4. Inicie os containers:
```bash
docker-compose up -d
```

## 🚀 Deployment

A aplicação pode ser implantada de duas formas:

### Usando Docker Compose
```bash
docker-compose up -d
```

### Usando a Imagem Docker
```bash
docker pull jeffemart/gotham:latest
docker run -p 8000:8000 jeffemart/gotham:latest
```

## 📦 CI/CD

O projeto utiliza GitHub Actions para:
- Build automático
- Testes
- Publicação da imagem Docker

Para usar o pipeline de CI/CD, configure os seguintes secrets no GitHub:
- `DOCKERHUB_USERNAME`: Seu usuário do Docker Hub
- `DOCKERHUB_TOKEN`: Token de acesso do Docker Hub

## 📝 Documentação

A documentação da API está disponível através do Swagger UI em:
```
http://localhost:8000/swagger/index.html
```

## 🔐 Segurança

- Autenticação via JWT
- Refresh tokens
- Senhas criptografadas com bcrypt
- CORS configurado
- Rate limiting
- Proteção contra ataques comuns

## 🤝 Contribuindo

1. Fork o projeto
2. Crie sua branch de feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 👤 Autor

Jefferson Martins
- LinkedIn: [jefferson-martins](https://www.linkedin.com/in/jefferson-martins-a6802b249/)
- Email: jefferson.developers@gmail.com

## 🙏 Agradecimentos

- Todos os contribuidores que ajudaram a tornar este projeto melhor
- A comunidade Go por ferramentas e bibliotecas incríveis

