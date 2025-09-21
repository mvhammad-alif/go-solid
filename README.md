# Go Solid - Clean Architecture Project

Go project skeleton using clean architecture.<br/>We use [echo](https://github.com/labstack/echo) as its framework and [wire](https://github.com/google/wire) as dependency injection tools.

## Features

- Clean Architecture with proper separation of concerns
- MySQL database with migration support
- Redis for caching (ready for implementation)
- Configuration management with Viper
- Dependency injection with Wire
- Docker support for MySQL and Redis

## Setup

### 1. Start Database Services

```bash
# Start MySQL and Redis using Docker
docker-compose up -d
```

### 2. Run the Application

```bash
go run server.go
```

### 3. Test the Endpoints

#### User Endpoints
- `GET http://localhost:1323/users/1` - Get user details

#### Post Endpoints
- `GET http://localhost:1323/sync` - Sync posts from external API
- `GET http://localhost:1323/items` - Get all posts

## Configuration

The application uses Viper for configuration management. You can configure the database and Redis settings in several ways:

### 1. Configuration File (config.yaml)
Create a `config.yaml` file in the root directory:

```yaml
database:
  host: "localhost"
  port: 3306
  user: "go_solid_user"
  password: "go_solid_pass"
  database: "go_solid"

redis:
  host: "localhost"
  port: 6379
```

### 2. Environment Variables
Set environment variables to override config values:

```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=go_solid_user
export DB_PASSWORD=go_solid_pass
export DB_NAME=go_solid
export REDIS_HOST=localhost
export REDIS_PORT=6379
```

### 3. Default Values
If no configuration is provided, the application will use these defaults:
- Database: localhost:3306 with user `go_solid_user` and database `go_solid`
- Redis: localhost:6379

## Database Schema

The application automatically creates the necessary tables on startup:

- `posts` - Stores post data with fields: id, user_id, title, body, timestamps
- `users` - Stores user data with fields: id, name, email, timestamps

## Project Structure

```
internal/
├── app/           # Application setup and dependency injection
├── config/        # Configuration management
├── database/      # Database migrations and utilities
├── delivery/      # HTTP handlers (presentation layer)
├── entity/        # Domain entities
├── repository/    # Data access layer
├── tools/         # Utility tools
└── usecase/       # Business logic layer
```
