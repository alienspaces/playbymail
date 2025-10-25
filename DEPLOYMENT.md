# Playbymail Deployment Guide (Heroku Buildpack, Monorepo, Automated Migrations)

This guide describes how to deploy the Playbymail full-stack application (Vue.js frontend + Go backend) to Heroku using a monorepo, GitLab CI/CD, and Heroku's Go buildpack. It details the configuration files and scripts that make the deployment work.

## Architecture Overview

- **Monorepo**: Both frontend (`frontend/`) and backend (`backend/`) code live in the same repository.
- **Frontend**: Vue.js (built to static files in `frontend/dist`)
- **Backend**: Go API server (serves API endpoints and static frontend files)
- **Database**: PostgreSQL (Heroku add-on)
- **Deployment**: Heroku Go buildpack, with automated DB migrations

## Key Configuration Files

### 1. `.gitlab-ci.yml` (Project Root)

Defines the CI/CD pipeline for testing, building, and deploying the app.

- **Stages**: `test`, `build`, `deploy`
- **Frontend**: Built in the `build-frontend` job, output saved as an artifact (`frontend/dist`)
- **Backend**: Tested in `test-backend`, deployed in `deploy-backend`
- **Deploy**: Uses `git subtree split` to push only the `backend/` directory to Heroku, so Heroku sees it as the app root.
- **Static Files**: Uses a `/tmp` artifact approach to copy the frontend build into the deploy branch after the split, ensuring static files are present without polluting the main branch.
- **Debugging**: The deploy job prints directory listings, file contents, and environment variables to help diagnose issues.

- For the most up-to-date deployment steps, see the `deploy-backend` job in `.gitlab-ci.yml`.

### 2. `backend/Procfile`

Tells Heroku how to start the app and run database migrations.

```Procfile
web: ls -l /app/public && ./bin/server
release: echo "=== DEBUG: PATH ===" && echo $PATH && echo "=== DEBUG: /app/bin ===" && ls -l /app/bin && echo "=== DEBUG: which migrate ===" && which migrate && echo "=== DEBUG: which river ===" && which river && /app/bin/migrate -verbose -path ./db -database "$DATABASE_URL" up && /app/bin/river migrate-up --database-url "$DATABASE_URL" up
```
- **web**: Starts the Go server built by the Go buildpack, with a debug step to list `/app/public`.
- **release**: Runs DB migrations using `migrate` and `river` before each deploy.

### 3. `backend/app.json`

Heroku app manifest for review apps and add-on provisioning.

```json
{
  "name": "playbymail-backend",
  "description": "Playbymail game backend API",
  "repository": "https://gitlab.com/alienspaces/playbymail",
  "env": {
    "PORT": {
      "description": "Port for the web server",
      "value": "8080"
    }
  },
  "formation": {
    "web": {
      "quantity": 1,
      "size": "basic"
    }
  },
  "addons": [
    { "plan": "heroku-postgresql:mini" }
  ],
  "buildpacks": [
    { "url": "https://github.com/heroku/heroku-buildpack-apt" },
    { "url": "heroku/go" }
  ]
}
```
- **Add-ons**: Provisions a Heroku Postgres database.
- **Buildpacks**: Uses the Apt buildpack to install OCR dependencies, then the official Go buildpack.

### 4. `backend/Aptfile`

Specifies system packages to install via the Apt buildpack.

```
tesseract-ocr
tesseract-ocr-dev
libtesseract-dev
libleptonica-dev
```
- **OCR Dependencies**: Installs Tesseract OCR and Leptonica libraries required for turn sheet scanning.

### 5. `backend/go.mod`

Defines Go module dependencies and the main module path.  
**Heroku uses this to detect and build your Go app.**

## CI/CD Pipeline Reference

The CI/CD pipeline is defined in the [`.gitlab-ci.yml`](.gitlab-ci.yml) file at the project root. This pipeline:
- Runs backend and frontend tests
- Builds the frontend and saves the output as an artifact
- Deploys the backend (with the built frontend) to Heroku using the `deploy-backend` job

For the most up-to-date and detailed steps, refer directly to the `deploy-backend` job in `.gitlab-ci.yml`.

## Preflight: Heroku Configuration Test

Before deploying, run the following script to verify your Heroku app, add-ons, and environment variables are set up correctly:

```sh
./tools/heroku-config-test
```

This script checks:
- Heroku CLI installation and authentication
- App accessibility
- Required environment variables (e.g., DATABASE_URL)
- PostgreSQL add-on presence
- Buildpacks, domains, dynos, and logs

**Requirements:**  
- [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli)  
- [jq](https://stedolan.github.io/jq/) (for JSON parsing)

## Deployment Flow

1. **Frontend Build**: GitLab CI builds the frontend and saves the output as an artifact (`frontend/dist`).
2. **Artifact Copy**: In the deploy job, the artifact is copied to `/tmp/frontend-dist` before the split.
3. **Subtree Deploy**: Only the `backend/` directory is pushed to Heroku, so Heroku sees it as the app root.
4. **Static Files**: After the split and checkout, the frontend build is copied from `/tmp/frontend-dist` to `public/` in the deploy branch and committed.
5. **Heroku Buildpack**: Heroku Go buildpack builds the Go app (`./cmd/server`), outputting `./bin/server`.
6. **Migration Binaries**: `migrate` and `river` binaries are built and committed in the deploy branch, so they are available for the release phase.
7. **Release Phase**: Heroku runs the `release` command from the `Procfile` to apply DB migrations.
8. **Web Process**: Heroku runs the Go server as the web process, serving static files from `/app/public`.

## How Database Migrations Are Run

### Overview

Database migrations are automatically applied during each deployment to Heroku, ensuring your database schema is always up to date with your application code.

### Migration Flow

1. **Migration Binaries Installed:**  
   The deploy job builds and commits the `migrate` and `river` binaries into the deploy branch, so they are present in `/app/bin` at runtime.
2. **Release Phase:**  
   The `release` process defined in `backend/Procfile` is run by Heroku before the new code is promoted to production. This process executes the following command:
   ```sh
   /app/bin/migrate -verbose -path ./db -database "$DATABASE_URL" up && /app/bin/river migrate-up --database-url "$DATABASE_URL" up
   ```
   - This applies all new migrations in the `db/` directory using `migrate`.
   - Then, it applies any River queue migrations.
3. **Failure Handling:**  
   If any part of the migration process fails, the deployment is aborted and the new code is **not** released. This ensures your app is never running with an out-of-sync schema.

### Key Files

- `backend/Procfile`  
  Defines the `release` process for running migrations.
- `backend/go.mod`  
  Go module definition.
- `backend/db/`  
  Contains migration files.

## Environment Variables

Set these on Heroku (via dashboard or CLI):

- `HEROKU_APP_NAME`: Your Heroku app name (for CI)
- `HEROKU_API_KEY`: Your Heroku API key (for CI)
- Any app-specific variables (e.g., `APP_HOME`, `ASSETS_PATH`, etc.)

The deployment and local scripts use `tools/environment` to validate and load required environment variables. For local development, see the `tools/start-backend` and `tools/start-frontend` scripts.

## Summary Table: Key Files

| File/Dir                | Purpose                                                      |
|-------------------------|--------------------------------------------------------------|
| `.gitlab-ci.yml`        | CI/CD pipeline, builds frontend, deploys backend to Heroku   |
| `backend/Procfile`      | Heroku process types: web and release (migrations)           |
| `backend/app.json`      | Heroku app manifest, add-ons, buildpack                      |
| `backend/go.mod`        | Go module definition                                         |
| `backend/public/`       | Built frontend assets (copied here by CI in deploy branch)   |

## Performance Considerations

- API responses are cached where appropriate
- Database connections are pooled

## Security

- CORS is configured for the same domain
- API endpoints are protected with authentication middleware
- Environment variables are used for sensitive configuration
- HTTPS is enforced by Heroku

## Monitoring

- Heroku provides basic monitoring and logging

If you need to update the deployment process, ensure all Heroku-related files (`Procfile`, `app.json`) are in the `backend/` directory, as this is what gets pushed to Heroku via subtree split.
