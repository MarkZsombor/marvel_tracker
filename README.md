# Marvel Champions Play Tracker

A web application for logging and tracking plays of the Marvel Champions board game. Built with Go, HTMX, and Tailwind CSS.

## Features

- Log Marvel Champions game sessions
- Track heroes, scenarios, and outcomes
- Server-side rendered HTML with HTMX for dynamic interactions
- Responsive design with Tailwind CSS
- SQLite database for simple deployment

## Technology Stack

- **Backend**: Go with Gin framework
- **Database**: SQLite
- **Frontend**: HTMX + Tailwind CSS
- **Testing**: Go built-in testing + testify

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd marvel_tracker
```

2. Install dependencies:

```bash
go mod tidy
```

3. Run the development server:

```bash
go run cmd/server/main.go
```

4. Open your browser to `http://localhost:8080`

### Development

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Build for production
go build -o bin/server cmd/server/main.go
```

## Project Structure

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
├── tests/              # Test files
├── plan.md             # Development plan
└── CLAUDE.md           # Claude AI context file
```

## Marvel Champions

Marvel Champions is a cooperative card game where players take on the roles of Marvel superheroes to defeat villains and save the world. This tracker helps log your game sessions including:

- Hero and aspect played
- Scenario and difficulty
- Game outcome and duration
- Notes and memorable moments

## Contributing

This is a personal learning project, but suggestions and feedback are welcome!

## License

MIT License - see LICENSE file for details
