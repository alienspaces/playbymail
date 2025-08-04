# Contributing to Playbymail

Thank you for your interest in contributing to the Playbymail project! This document provides guidance for local development, troubleshooting, and manual operations.

## Project Structure Overview

| Component | Location | Description |
|-----------|----------|-------------|
| Git Repository | `./playbymail` | Main project repository |
| Frontend UI | `./playbymail/frontend` | User interface components |
| Backend Server | `./playbymail/backend` | Server-side (Golang) |

### Backend Directory Structure

| Directory | Location | Purpose |
|-----------|----------|---------|
| Main Executables | `./playbymail/backend/cmd` | Main executables |
| Core Library | `./playbymail/backend/core` | Core library functions |
| Database Migrations | `./playbymail/backend/db` | Database migrations files |
| API Schemas | `./playbymail/backend/schema` | JSON schema definitions for API endpoints |
| Business Logic | `./playbymail/backend/internal/domain` | Business rules |
| Test Harness | `./playbymail/backend/internal/harness` | Test data setup and teardown |
| Database Records | `./playbymail/backend/internal/record` | Database record definitions |
| Repository Layer | `./playbymail/backend/internal/repository` | Database repository implementations |
| CLI Runner | `./playbymail/backend/internal/runner/cli` | CLI executable implementation |
| Server Runner | `./playbymail/backend/internal/runner/server` | Server executable implementation |

### Local Development Scripts

The `tools/` directory contains scripts to help you set up and run the project locally. Here are the most relevant scripts for contributors:

- **Start Everything:**

  ```sh
  ./tools/start
  ```

  Starts the backend API server, frontend development server, and ensures the database is running. This is the recommended way to start all services for local development.  
  - Backend API: http://localhost:8080  
  - Frontend: http://localhost:3000

- **Stop Everything:**

  ```sh
  ./tools/stop
  ```

  Stops the backend, frontend, and database. Uses PID files for safe process management.

- **Start Backend Only:**

  ```sh
  ./tools/start-backend
  ```

  Builds and starts the Go backend API server. Ensures the database is running.

- **Start Frontend Only:**

  ```sh
  ./tools/start-frontend
  ```

  Starts the Vue.js frontend development server.

- **Stop Backend Only:**

  ```sh
  ./tools/stop-backend
  ```

  Stops the backend server process (using PID file).

- **Stop Frontend Only:**

  ```sh
  ./tools/stop-frontend
  ```

  Stops the frontend development server (using PID file).

- **Database Setup:**

  ```sh
  ./tools/db-setup
  ```

  Starts the local database (in Docker), runs migrations, and loads test/reference data.

- **Database Setup for Tests:**

  ```sh
  ./tools/db-setup-test
  ```

  Prepares the database with test data for running backend tests.

- **Run All Tests:**

  ```sh
  ./tools/test-all
  ```

  Runs all backend and frontend tests.

- **Heroku Configuration Test:**

  ```sh
  ./tools/heroku-config-test
  ```

  Checks that your Heroku app is configured correctly for deployment.

- **Other Utilities:**
  - **Database Migration:**  
    `./tools/db-migrate-up`, `./tools/db-migrate-down`, `./tools/db-migrate-create`
  - **Database Connect:**  
    `./tools/db-connect`, `./tools/db-connect-qa`
  - **Data Loaders:**  
    `./tools/db-load-test-data`, `./tools/db-load-test-reference-data`
  - **Export Table Data:**  
    `./tools/db-export-table-data`
  - **Test Helpers:**  
    `./tools/test-backend`, `./tools/test-backend-ci`, `./tools/test-backend-core`, `./tools/test-backend-internal`, `./tools/test-frontend`
  - **Build Tools:**  
    `./tools/build-frontend`
  - **Database Management:**  
    `./tools/db-start`, `./tools/db-stop`
  - **Development Utilities:**  
    `./tools/retry`

Refer to the comments at the top of each script in `tools/` for more details and additional utilities.
