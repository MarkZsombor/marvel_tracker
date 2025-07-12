# Marvel Champions Play Tracker - Development Plan

## Project Overview

A web application to log and track plays for the Marvel Champions board game, built with HTMX, Tailwind CSS, Go, and SQLite.

## Technology Stack

- **Frontend**: HTMX + Tailwind CSS with server-side rendered HTML
- **Backend**: Go with Gin framework
- **Database**: SQLite with database/sql
- **Testing**: Go built-in testing + testify
- **CI/CD**: GitHub Actions

## Phase 1: Project Foundation

### 1. Initial Go Setup

- [ ] Initialize Go module (`go mod init marvel_tracker`)
- [ ] Set up project directory structure
- [ ] Install core dependencies:
  - `github.com/gin-gonic/gin` - Web framework
  - `github.com/mattn/go-sqlite3` - SQLite driver
  - `github.com/stretchr/testify` - Testing assertions
- [ ] Create basic main.go with Gin server
- [ ] Set up basic routing structure

### 2. Project Structure

```
marvel_tracker/
├── cmd/server/          # Main application entry point
│   └── main.go
├── internal/
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Data models and database logic
│   ├── middleware/      # Custom middleware
│   └── config/          # Configuration management
├── templates/           # HTML templates
├── static/             # Static assets (CSS, JS, images)
├── migrations/         # Database migration files
├── tests/              # Test files
├── go.mod
├── go.sum
└── README.md
```

### 3. Database Setup

- [ ] Create SQLite database initialization
- [ ] Design initial schema using a relational model:
  - **`heroes`** (id, name) - _Master list of heroes._
  - **`scenarios`** (id, name) - _Master list of scenarios._
  - **`plays`** (id, date, outcome, notes, scenario_id, difficulty) - _Records a single game session, linking to one scenario._
  - **`decks`** (id, play_id, hero_id, aspect) - _Links a play to the heroes used, storing play-specific data like the aspect._
- [ ] Plan for seeding initial `heroes` and `scenarios` data (e.g., via migration).
- [ ] Set up database connection and basic CRUD operations for the models.
- [ ] Create a simple migration system.

### 4. Basic Web Server

- [ ] Set up Gin router with basic routes
- [ ] Configure HTML template rendering
- [ ] Create basic layout template
- [ ] Set up static file serving
- [ ] Implement basic error handling

## Phase 2: Frontend Integration

### 5. HTMX Integration

- [ ] Add HTMX to templates via CDN
- [ ] Create first HTMX endpoint (simple form submission)
- [ ] Set up partial template rendering for HTMX responses
- [ ] Implement basic navigation without page refreshes

### 6. Tailwind CSS Setup

- [ ] Add Tailwind CSS via CDN for initial development.
- [ ] Set up Tailwind CLI build process for customization and production optimization.
- [ ] Create base layout with Tailwind classes.
- [ ] Design responsive navigation.
- [ ] Plan for handling UI loading states and form validation feedback.

### 7. Basic UI Components

- [ ] Header/navigation component
- [ ] Form components (input, button, select)
- [ ] Table component for displaying plays
- [ ] Modal/dialog component

## Phase 3: Core Functionality

### 8. Play Logging Features

- [ ] Create "New Play" form with HTMX
- [ ] Implement play creation handler
- [ ] Display list of plays with sorting/filtering
- [ ] Basic play editing functionality
- [ ] Play deletion with confirmation

### 9. Marvel Champions Specific Features

- [ ] Hero selection dropdown
- [ ] Scenario selection dropdown
- [ ] Difficulty tracking
- [ ] Outcome recording (win/loss)

## Phase 4: Testing & Quality

### 10. Testing Setup

- [ ] Set up Go testing structure
- [ ] Unit tests for models and handlers
- [ ] Integration tests for API endpoints
- [ ] Test database setup (separate test DB)
- [ ] Basic test coverage reporting

### 11. Code Quality

- [ ] Set up Go formatting and linting
- [ ] Add input validation
- [ ] Implement proper error handling
- [ ] Add logging throughout application

## Phase 5: Deployment Preparation

### 12. GitHub Actions CI/CD

- [ ] Create workflow for running tests
- [ ] Add linting and formatting checks
- [ ] Set up build process
- [ ] Basic deployment pipeline structure

### 13. Production Readiness

- [ ] Environment configuration
- [ ] Database migration strategy
- [ ] Static asset optimization
- [ ] Basic security headers
- [ ] Health check endpoint

## Phase 6: Future Enhancements (Post-MVP)

- [ ] Detailed play statistics and analytics
- [ ] Hero/scenario win rate tracking
- [ ] Play session photos/notes
- [ ] Import/export functionality
- [ ] User authentication (if multi-user needed)
- [ ] Advanced filtering and search
- [ ] Mobile-responsive improvements

## Learning Resources

- **Go**: Official Go documentation, Go by Example
- **HTMX**: HTMX documentation, HTMX examples
- **Gin**: Gin framework documentation
- **SQLite**: SQLite tutorial, SQL basics

## Getting Started Commands

```bash
# Initialize project
go mod init marvel_tracker

# Install dependencies
go get github.com/gin-gonic/gin
go get github.com/mattn/go-sqlite3
go get github.com/stretchr/testify

# Run development server
go run cmd/server/main.go

# Run tests
go test ./...

# Format code
go fmt ./...
```

## Notes

- Start with simple functionality and iterate
- Focus on learning Go and HTMX patterns before adding complexity
- Keep database schema simple initially
- Use HTMX progressively - start with simple form submissions
- Tailwind can be added via CDN initially, optimize later
