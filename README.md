# Go Solid - Clean Architecture Project

Go project skeleton using clean architecture.<br/>We use [echo](https://github.com/labstack/echo) as its framework and [wire](https://github.com/google/wire) as dependency injection tools.

## Features

- Clean Architecture with proper separation of concerns
- MySQL database with migration support
- **Redis caching with 1-minute TTL** for improved performance
- Configuration management with Viper
- Dependency injection with Wire
- Docker support for MySQL and Redis
- **Resilient external API calls with exponential backoff** using cenkalti/backoff

## Setup

### 1. Start Database Services

```bash
# Start MySQL and Redis using Docker
make docker-up
```

### 2. Run Database Migrations

```bash
# Run database migrations to create tables
make migration
```

### 3. Run the Application

The application now has two separate entry points:

#### HTTP Server (API Server)
```bash
# Run the HTTP server with all endpoints
make server
```

#### Cron Service (Background Jobs)
```bash
# Run only the cron service for background jobs
make cron
```

### 4. Test the Endpoints

#### User Endpoints
- `GET http://localhost:1323/users/1` - Get user details

#### Post Endpoints
- `GET http://localhost:1323/sync` - Sync posts from external API
- `GET http://localhost:1323/items` - Get all posts

### 5. Running Both Services

You can run both services simultaneously in different terminals:

**Terminal 1 - HTTP Server:**
```bash
make server
```

**Terminal 2 - Cron Service:**
```bash
make cron
```

Or use a process manager like supervisor or systemd to run them as services.

## Makefile Commands

The project includes a comprehensive Makefile with convenient commands:

### Development Commands
- `make help` - Show all available commands
- `make server` - Run the HTTP server
- `make cron` - Run the cron service
- `make migration` - Run database migrations

### Build Commands
- `make build` - Build all binaries (server, cron, migration)
- `make clean` - Clean build artifacts
- `make prod-build` - Build optimized binaries for production

### Database Commands
- `make docker-up` - Start MySQL and Redis containers
- `make docker-down` - Stop MySQL and Redis containers

### Testing Commands
- `make test` - Run tests

### Quick Development
- `make dev` - Build and run the server in development mode

### Example Usage
```bash
# Quick start - setup everything and run server
make docker-up
make migration
make server

# Or run cron service
make cron

# Build for production
make prod-build
```

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

## Redis Caching

The application implements intelligent caching for improved performance:

### Caching Strategy
- **Cache First**: Check Redis cache before hitting the database
- **Cache Miss**: Fetch from MySQL and store result in Redis
- **TTL**: 1-minute expiration time for cached data
- **Cache Key**: `posts:all` for the main posts listing

### Architecture
- **Tools Layer** (`internal/tools/redis.go`): Redis client wrapper with clean methods
- **Repository Layer** (`internal/repository/post/post.go`): Implements caching logic using Redis tools
- **Usecase Layer** (`internal/usecase/post/post.go`): Calls repository with caching enabled

### Redis Methods Available
- `Fetch(ctx, key)` - Get data from Redis
- `Store(ctx, key, value, expiration)` - Store data in Redis with TTL
- `Delete(ctx, key)` - Remove data from Redis
- `Exists(ctx, key)` - Check if key exists in Redis

### Performance Benefits
- **Reduced Database Load**: Frequent requests served from Redis
- **Faster Response Times**: In-memory Redis operations are much faster than MySQL queries
- **Automatic Expiration**: Stale data is automatically cleaned up after 1 minute
- **Fallback Strategy**: If Redis is unavailable, falls back to direct database queries

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

## External API Resilience

The `SyncPosts` functionality includes robust error handling and retry logic:

### Backoff Strategy
- **Exponential backoff** with jitter to avoid thundering herd problems
- **Initial interval**: 1 second
- **Maximum interval**: 30 seconds
- **Maximum elapsed time**: 5 minutes
- **Multiplier**: 2x (each retry doubles the wait time)

### Retry Logic
- ✅ **Retries on network errors** (connection failures, timeouts)
- ✅ **Retries on server errors** (5xx status codes)
- ❌ **No retry on client errors** (4xx status codes) - these are permanent
- ❌ **No retry on unexpected status codes** - these indicate API changes
- **Context-aware**: Respects request context and cancels on timeout

### Error Handling
- Individual post creation failures don't stop the entire sync process
- Failed posts are logged with warnings but processing continues
- Comprehensive error messages for debugging

## Cron Service

The application includes a built-in cron service that runs scheduled tasks:

### Architecture
- **Delivery Layer**: Contains `CronHandler` with job logic (moved from tools layer)
- **Provider Layer**: Defines cron jobs and their schedules in `internal/app/provider.go`
- **Tools Layer**: Generic `CronService` that executes scheduled jobs
- **Dependency Injection**: Jobs are wired together using Google's Wire

### Scheduled Jobs
- **Sync Posts**: Runs every 15 minutes to fetch new posts from the external API
- **Configurable**: Easy to add more scheduled jobs by modifying `provideCronJobs()` in provider.go

### Running the Cron Service
```bash
# Run only the cron service (no HTTP server)
go run cmd/cron/main.go
```

### Cron Service Features
- **Clean Architecture**: Proper separation of concerns with jobs defined in delivery layer
- **Graceful startup and shutdown**: Properly handles signals for clean termination
- **Detailed logging**: Logs all cron job executions with timestamps and status
- **Error handling**: Continues running even if individual jobs fail
- **Configurable timeouts**: Each job has a timeout to prevent hanging
- **Independent operation**: Can run without the HTTP server

### Cron Job Details
- **Schedule**: `*/15 * * * *` (every 15 minutes)
- **Timeout**: 10 minutes per execution
- **Logging**: All executions are logged with success/failure status
- **Error recovery**: Failed jobs don't stop the service

### Adding New Cron Jobs
To add a new cron job, modify the `provideCronJobs()` function in `internal/app/provider.go`:

```go
func provideCronJobs(cronHandler *delivery.CronHandler) []tools.CronJob {
    return []tools.CronJob{
        {
            Name:     "sync_posts",
            Schedule: "*/15 * * * *", // Every 15 minutes
            Job:      cronHandler.SyncPostsJob,
        },
        {
            Name:     "new_job",
            Schedule: "0 2 * * *", // Daily at 2 AM
            Job:      cronHandler.NewJob,
        },
    }
}
```
