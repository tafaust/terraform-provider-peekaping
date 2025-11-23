# Terraform Provider Tests

This directory contains tests for the Terraform provider.

## Running Tests

### With Docker Compose (Recommended)

The easiest way to run tests is using Docker Compose, which automatically starts a Peekaping API server:

```bash
make test-native-docker
```

This will:
1. Start a Peekaping API server using Docker Compose
2. Wait for the server to be ready
3. Run the Terraform tests
4. Stop the Docker Compose services

### Manual Setup

If you prefer to run tests against an existing API server:

1. Set environment variables:
   ```bash
   export PEEKAPING_ENDPOINT=http://localhost:8034
   export PEEKAPING_EMAIL=your-email@example.com
   export PEEKAPING_PASSWORD=your-password
   # OR use API key:
   export PEEKAPING_API_KEY=your-api-key
   ```

2. Run tests:
   ```bash
   make test-native
   ```

## Test Files

- `main.tf` - Main test configuration
- `variables.tf` - Test variables
- `monitor.tftest.hcl` - Monitor resource tests
- `deferred_data_sources.tf` - Tests for deferred data sources (fixes for "inconsistent final plan" errors)
- `deferred_data_sources.tftest.hcl` - Test cases for deferred data sources

## Docker Compose Test Setup

The test setup uses `docker-compose.test.yml` which:
- Uses the SQLite bundle from `peekaping-upstream` submodule
- Maps port 8034 (expected by tests) to 8383 (container port)
- Stores test data in `.test-data/` directory (gitignored)
- Includes health checks to ensure the API is ready before tests run

## Troubleshooting

### API Server Not Starting

If the API server doesn't start:
```bash
# Check logs
docker compose -f docker-compose.test.yml logs peekaping

# Stop and clean up
make test-docker-down
```

### Port Already in Use

If port 8034 is already in use, you can modify `docker-compose.test.yml` to use a different port and update `PEEKAPING_ENDPOINT` accordingly.

### Submodule Not Initialized

If you get errors about missing `peekaping-upstream`:
```bash
git submodule update --init --recursive
```

