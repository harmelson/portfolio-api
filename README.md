# Portfolio API

A portfolio-ready REST API for managing users, subscription plans, and user subscriptions. It uses Go, PostgreSQL, `pgx`, Google ID token authentication, and a layered architecture with handlers, services, and repositories.

## Features

- Google ID token authentication
- User and subscription management
- PostgreSQL transactions with `pgx`
- Request body limits and authentication middleware
- Graceful HTTP server shutdown

## Run with Docker

From the repository root, start the API and PostgreSQL database:

```bash
docker compose up --build
```

The API is available at `http://localhost:8080` and PostgreSQL is exposed on port `5432`.

The database is initialized with:

- Plans: `free`, `basic`, and `plus`
- Test user: `dev@example.com` (`google_id: dev-test-user`)

To stop the services, press `Ctrl+C`. To remove containers and database data, run:

```bash
docker compose down -v
```

## Endpoints

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| `GET` | `/health` | Returns the API health status. | None |
| `GET` | `/ready` | Returns the API readiness status. | None |
| `GET` | `/api/v1/users/get` | Returns the authenticated user. | Google ID token |
| `POST` | `/api/v1/users` | Creates the authenticated user. | Google ID token |
| `GET` | `/api/v1/subscriptions/me` | Returns the authenticated user's subscription and plan. | Google ID token or development token |

## Authentication

Send a Google ID token in the `x-tiger-auth` header:

```bash
curl -H "x-tiger-auth: Bearer <google-id-token>" http://localhost:8080/api/v1/users/get
```

When `DEV_AUTH_ENABLED=true`, the subscription endpoint also accepts the development header configured by `DEV_AUTH_TOKEN`:

```bash
curl -H "x-dev-auth: dev-token" http://localhost:8080/api/v1/subscriptions/me
```

## Environment Variables

Use [`.env.example`](.env.example) as a reference for the available configuration variables.
