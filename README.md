# GoBank

A simple banking simulation API built with Go, featuring user authentication, account management, and money transfers.

## Features

- User registration and login with JWT authentication
- Account creation and balance management
- Secure money transfers between accounts
- Asynchronous processing with RabbitMQ
- PostgreSQL database with GORM ORM
- RESTful API with Gin framework

## Prerequisites

- Go 1.19+
- PostgreSQL
- RabbitMQ
- Docker (optional, for containerized setup)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/trueMNOX/GoBank.git
   cd GoBank
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Set up environment variables (create a `.env` file or set them directly):
   - `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
   - `JWT_SECRET`
   - `RABBITMQ_URL`
   - `SERVER_PORT`

4. Run database migrations:

   ```bash
   go run migrations/migrate.go
   ```

## Usage

### Running the API Server

```bash
go run cmd/api/main.go
```

The server will start on the configured port (default: 8080).

### Running the Worker

```bash
go run cmd/worker/main.go
```

### API Endpoints

- `POST /api/register` - Register a new user
- `POST /api/login` - Login user
- `POST /api/accounts` - Create account (authenticated)
- `GET /api/accounts` - List user accounts (authenticated)
- `POST /api/transfers` - Transfer money (authenticated)

## Docker

To run with Docker Compose:

```bash
docker-compose up --build
```

## License

This project is licensed under the MIT License.
