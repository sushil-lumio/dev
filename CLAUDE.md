# CLAUDE.md — [Project Name]

> Project-specific guide for Claude Code. See the [team CLAUDE.md](link-to-global) for org-wide conventions.

## Overview

<!-- One-line description of what this service/app does -->

## Quick Start

```bash
# Install dependencies
[command]

# Run locally
[command]

# Run tests
[command]

# Build
[command]

# Deploy (specify environment)
[command]
```

## Project Structure

```
[root]/
├── [key directories and files]
└── [explain what lives where]
```

## Key Dependencies

<!-- List the critical libraries/frameworks and what they're used for -->
<!-- Example:
- `aws-lambda-go` — Lambda handler + API Gateway event types
- `go-sql-driver/mysql` — MySQL driver
-->

## Environment Variables

<!-- List required env vars. Reference .env.example if it exists -->
<!-- Example:
| Variable | Description | Required |
|----------|-------------|----------|
| DB_ENDPOINT | MySQL RDS endpoint | Yes |
| DB_USER | Database username | Yes |
-->

## Database

<!-- Which database(s) does this repo use? MySQL, DynamoDB, Redis, etc. -->
<!-- Connection patterns, important tables/collections -->

## API Endpoints

<!-- If this is an API service, list the endpoints -->
<!-- Example:
| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
-->

## Deployment

<!-- How is this deployed? Which environment? -->
<!-- Example:
- Platform: AWS Lambda (ap-south-1)
- Terraform: `terraform/` directory
- Deploy dev: `make deploy TF_VAR_FILE=dev.tfvars TF_STATE_FILE=terraform_dev.tfstate`
-->

## Architecture Decisions

<!-- Any non-obvious design choices and WHY they were made -->

## Gotchas

<!-- Repo-specific pitfalls that catch people off guard -->
<!-- Example:
- The `/dev` and `/prod` path prefixes must be stripped in main.go routing
- Static JSON files must be added to Makefile zip target manually
-->

## Testing

<!-- How to run tests, what's covered, what's not -->
<!-- If no tests exist yet, note that and what should be tested first -->
