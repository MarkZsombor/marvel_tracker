# Marvel Champions Play Tracker - Claude Context

## Project Overview

A web application for logging and tracking plays of the Marvel Champions board game. Built as a learning project to explore HTMX and Go development.

## Technology Stack

- **Backend**: Go with Gin framework
- **Database**: SQLite with database/sql package
- **Frontend**: Server-side rendered HTML with HTMX and Tailwind CSS
- **Testing**: Go built-in testing + testify assertions
- **CI/CD**: GitHub Actions

## Development Environment

### Commands

```bash
# Development server
go run cmd/server/main.go

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Install dependencies
go mod tidy

# Build for production
go build -o bin/server cmd/server/main.go
```

### Project Structure

```
marvel_tracker/
├── cmd/server/          # Main application entry point
├── internal/
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Data models and database logic
│   ├── middleware/      # Custom middleware
│   └── config/          # Configuration management
├── templates/           # HTML templates
├── static/             # Static assets (CSS, JS, images)
├── migrations/         # Database migration files
└── tests/              # Test files
```

## Coding Conventions

### Go Style

- Follow standard Go conventions (gofmt, golint)
- Use receiver names that are 1-2 characters and consistent
- Package names should be lowercase, single words
- Interface names should end with -er when possible
- Use table-driven tests for multiple test cases

### File Organization

- Group related functionality in packages under `internal/`
- Keep handlers thin - business logic should be in models or services
- Use dependency injection for database connections
- Templates should be organized by feature/page

### Database

- Use snake_case for table and column names
- Include created_at and updated_at timestamps on main entities
- Use migrations for schema changes
- Keep database queries in model files, not handlers

### HTMX Patterns

- Use `hx-` attributes for dynamic behavior
- Return partial HTML templates for HTMX requests
- Use `hx-target` to specify where responses should be rendered
- Implement proper error handling for HTMX requests

### HTML Templates

- Use semantic HTML elements
- Apply Tailwind classes for styling
- Keep templates DRY with partials and layouts
- Use template functions for common formatting

## Key Dependencies

```go
// Core dependencies
github.com/gin-gonic/gin           // Web framework
github.com/mattn/go-sqlite3        // SQLite driver

// Testing
github.com/stretchr/testify        // Test assertions

// Potential future additions
// github.com/golang-migrate/migrate // Database migrations
// github.com/go-playground/validator // Input validation
```

## Marvel Champions Domain

### Core Entities

- **Play**: A single game session (date, duration, outcome, notes)
- **Hero**: Marvel character being played (name, aspect color)
- **Scenario**: Villain/situation being fought (name, difficulty)

### Game Concepts

- **Aspects**: Leadership (yellow), Justice (blue), Aggression (red), Protection (green)
- **Difficulty**: Standard I, Standard II, Expert I, Expert II, Heroic I-IV
- **Outcomes**: Win, Loss (with optional win condition details)

## Development Notes

- Start simple - basic CRUD for plays before adding complexity
- HTMX should enhance the experience, not be required for basic functionality
- Database schema should accommodate future features but stay simple initially
- Focus on learning Go idioms and HTMX patterns
- Test coverage should be reasonable but not obsessive in early development

## Deployment Considerations

- SQLite database file location and permissions
- Static file serving configuration
- Environment variable configuration
- Graceful server shutdown handling
- Health check endpoint for monitoring

## Learning Resources

- [Go by Example](https://gobyexample.com/)
- [HTMX Documentation](https://htmx.org/docs/)
- [Gin Framework Guide](https://gin-gonic.com/docs/)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
