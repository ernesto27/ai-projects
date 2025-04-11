# Task Manager

A web-based task management application built with Go, Gin, and PostgreSQL.

## Features

- User authentication and authorization
- Task creation and management
- Web-based interface with responsive design
- RESTful API endpoints

## Prerequisites

- [Go](https://golang.org/dl/) (version 1.24.0 or later)
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/downloads)

## Local Setup

### 1. Clone the repository

```bash
git clone https://github.com/ernesto/task-manager.git
cd task-manager
```

### 2. Set up environment variables

Create a `.env` file in the root directory of the project:

```bash
# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=taskmanager
DB_SSLMODE=disable

# Application settings
PORT=8080
JWT_SECRET=your_jwt_secret_key
```

### 3. Start the PostgreSQL database using Docker

```bash
docker-compose up -d
```

This will start a PostgreSQL server on port 5432.


### 4. Build and run the application

```bash
go run main.go
```

Or run without building:

```bash
go run main.go
```

### 5. Access the application

Open your browser and navigate to:

```
http://localhost:8080
```

## API Endpoints

- `GET /api/tasks` - Get all tasks
- `POST /api/tasks` - Create a new task
- `GET /api/tasks/:id` - Get a specific task
- `PUT /api/tasks/:id` - Update a task
- `DELETE /api/tasks/:id` - Delete a task
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Log in and get JWT token

## Development

## Troubleshooting

### Database Connection Issues

If you have problems connecting to the database, make sure:
1. Docker containers are running (`docker-compose ps`)
2. Environment variables match your Docker Compose configuration
3. Database is healthy (`docker-compose logs postgres`)

### Authentication Issues

Make sure the JWT_SECRET is properly set in your .env file and that you're passing the token correctly in the Authorization header.

## License

[MIT](LICENSE)