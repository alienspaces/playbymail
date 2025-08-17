# PlayByMail

A platform to revitalize play-by-mail (PBM) games for children, teenagers, and young adults, supporting a wide variety of genres through **physical printed materials** mailed to players.

## What is PlayByMail?

PlayByMail is a **physical play-by-mail gaming platform** that:

- **Generates printed game materials** (maps, forms, puzzles, story content) to be mailed to players
- **Processes returned materials** from players via prepaid envelopes
- **Automates turn processing** based on returned materials with minimal human intervention
- **Supports multiple game genres** including RPGs, strategy games, adventure games, and more
- **Focuses on age-appropriate content** with parental controls and moderation

## How It Works

1. **Game Design**: Game designers create games using the Studio tools
2. **Game Management**: Game managers launch game instances and manage players
3. **Print Generation**: The system generates personalized print packs for each player
4. **Physical Mailing**: Materials are printed and mailed to players on a schedule (weekly/fortnightly)
5. **Player Response**: Players complete forms, solve puzzles, and return materials via prepaid envelopes
6. **Automated Processing**: The system scans and processes returned materials to update game state
7. **Next Turn**: New print packs are generated and mailed for the next turn

## Game Genres Supported

- **Role-Playing Games (RPGs)**: Character-driven adventures and storytelling
- **Wargames & Strategy**: Military campaigns, diplomacy, and resource management
- **Adventure & Exploration**: Dungeon crawls, space exploration, treasure hunting
- **Economic Games**: Empire building, trading, and resource management
- **Sports Management**: Team management and strategic decision-making
- **Science Fiction**: Space opera, colonization, and interstellar conflict
- **Fantasy Kingdom Management**: Medieval fantasy world building and conquest
- **Mystery & Detective**: Puzzle-solving and crime investigation
- **Collaborative Storytelling**: Group narrative building and shared adventures
- **Board Game Adaptations**: Classic games adapted for PBM format

## Platform Architecture

### 1. **Game Designer Studio** (`/studio`)
Tools for creating and managing play-by-mail games:
- Game design and creation tools
- Storyline and narrative development
- Map and visual content creation
- Form and puzzle design
- Game rule configuration
- Content testing and preview
- Game publishing and version management

### 2. **Game Management** (`/management`)
Tools for operating and administering running games:
- Subscribe to or purchase games to run
- Configure game runtime parameters
- Launch and manage game instances
- Process player turns and game state updates
- Review running game states and player progress
- Handle player communications and support
- Generate and manage print/mail operations

### 3. **Platform Administration** (`/admin`)
Platform-level administration and business operations:
- User account management and moderation
- Customer support and ticket management
- Platform-wide analytics and reporting
- Billing and subscription management
- System configuration and maintenance
- Content moderation and policy enforcement
- Platform announcements and communications

## Key Features

- **Physical Print & Mail**: All game content delivered via printed materials
- **Automated Turn Processing**: Minimal human intervention required
- **Personalized Content**: Each player receives customized materials
- **Return Processing**: Automated scanning and processing of returned materials
- **Age-Appropriate Content**: Built-in moderation and parental controls
- **Privacy-First**: Complete data removal available upon request
- **Accessibility**: Designed for players with disabilities
- **AI Integration**: AI-generated content with human review safeguards

## Tech Stack

- **Frontend**: Vue.js with Vite (for admin/management interfaces)
- **Backend**: Go with Gin framework (game engine and processing)
- **Database**: PostgreSQL with Row Level Security (RLS)
- **Print Generation**: PDF generation and print-ready output
- **Testing**: Playwright for end-to-end testing, Vitest for unit testing
- **Deployment**: Heroku

## Quick Start

### Prerequisites

- Node.js 18+ (use `nvm` if available)
- Go 1.21+
- PostgreSQL 15+
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd playbymail
   ```

2. **Install dependencies**
   ```bash
   # Install Node.js dependencies
   npm install
   
   # Install Go dependencies
   cd backend
   go mod download
   cd ..
   ```

3. **Environment setup**
   ```bash
   # Copy environment template
   cp .env.example .env
   
   # Edit .env with your configuration
   # Database connection, API keys, etc.
   ```

4. **Database setup**
   ```bash
   # Start database
   ./tools/db-start
   
   # Run migrations and load test data
   ./tools/db-setup
   ```

### Running the Application

#### Start Everything (Recommended for development)
```bash
# Start backend, frontend, and database
./tools/start
```

#### Start Individual Services
```bash
# Start only backend (includes database setup)
./tools/start-backend

# Start only frontend
./tools/start-frontend

# Start only database
./tools/db-start
```

#### Stop Services
```bash
# Stop all services
./tools/stop

# Stop specific services
./tools/stop-backend
./tools/stop-frontend
./tools/db-stop
```

### Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Database**: localhost:5432 (if using local PostgreSQL)

## Development

### Project Structure

```
playbymail/
├── frontend/              # Vue.js admin/management interfaces
├── backend/               # Go backend (game engine + processing)
├── tools/                 # Development and deployment scripts
├── playwright/            # End-to-end tests
├── .env                   # Environment configuration
├── .env.example           # Environment template
└── README.md              # This file
```

### Available Tools

The `./tools/` directory contains scripts for common development tasks:

#### Application Management
- `./tools/start` - Start complete application stack
- `./tools/stop` - Stop all services
- `./tools/start-backend` - Start backend with database
- `./tools/start-frontend` - Start frontend development server

#### Database Management
- `./tools/db-setup` - Complete database setup (migrations + test data)
- `./tools/db-start` - Start database service
- `./tools/db-stop` - Stop database service
- `./tools/db-connect` - Connect to local database
- `./tools/db-query` - Execute SQL queries
- `./tools/db-migrate-create [name]` - Create new migration

#### Testing
- `./tools/test-all` - Run all tests (frontend + backend)
- `./tools/test-frontend` - Run frontend tests
- `./tools/test-backend` - Run backend tests
- `./tools/test-playwright` - Run end-to-end tests

### Database Development

```bash
# Create new migration
./tools/db-migrate-create add_new_feature

# Apply migrations
./tools/db-migrate-up

# Rollback migrations
./tools/db-migrate-down

# Load test data
./tools/db-load-test-data
```

### Frontend Development

```bash
# Start development server
./tools/start-frontend

# Build for production
./tools/build-frontend

# Run tests
./tools/test-frontend
```

### Backend Development

```bash
# Start backend server
./tools/start-backend

# Run tests
./tools/test-backend

# Run specific test categories
./tools/test-backend-core
./tools/test-backend-internal
```

## Testing

### Running Tests

```bash
# Run all tests
./tools/test-all

# Run specific test suites
./tools/test-frontend      # Frontend unit tests
./tools/test-backend       # Backend tests
./tools/test-playwright    # End-to-end tests
```

### Playwright End-to-End Testing

For comprehensive UI testing, see the [Playwright documentation](playwright/README.md).

```bash
# Prerequisites: Backend running, database setup, frontend built
./tools/start-backend
./tools/build-frontend
./tools/test-playwright
```

## Deployment

### Heroku Deployment

```bash
# Set Heroku app name
heroku git:remote -a playbymail

# Deploy
git push heroku main

# Check logs
heroku logs --tail
```

### Environment Configuration

```bash
# Set environment variables
heroku config:set DATABASE_URL=your_database_url
heroku config:set JWT_SECRET=your_jwt_secret

# View current config
heroku config
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests to ensure everything works
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Guidelines

- Follow existing code patterns and conventions
- Write tests for new functionality
- Update documentation as needed
- Use the provided tools for common tasks
- Keep commits focused and descriptive

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 3000, 8080, and 5432 are available
2. **Database connection**: Check `.env` configuration and database status
3. **Node version**: Use `nvm` to ensure correct Node.js version
4. **Go modules**: Run `go mod download` in backend directory

### Getting Help

- Check the logs: `./tools/start-backend` shows backend logs
- Review environment configuration in `.env`
- Check database status: `./tools/db-connect`
- Run tests to identify issues: `./tools/test-all`

## License

[ISC License](LICENSE)

## Support

For issues and questions:
- Create an issue in the repository
- Check existing documentation
- Review test output for debugging information

---

**Important**: PlayByMail is a **physical play-by-mail gaming platform**, not a web-based game. The web interface is for game designers, managers, and administrators to create and manage games that are physically mailed to players.
