# Docker Startup Fix - Implementation Notes

## Problem
The core service experienced "pending requests" after starting with `docker-compose up`. Requests would hang until manually restarting the container. This was caused by:

1. **Race condition**: Core service started before all dependencies were ready
2. **Missing healthcheck**: Docker marked the service as "started" before it could accept connections
3. **No startup coordination**: Database migrations and RabbitMQ connection had no retry/wait logic

## Solution Implemented

### 1. Wait-for-it Script (`wait-for-it.sh`)
- Created a lightweight shell script that waits for services to be available
- Uses `netcat` to check port availability before starting the main application
- Configurable timeout (60s default)
- Waits for database on `db:5432` before starting the Go server

### 2. Updated Dockerfile (`Dockerfile.yml`)
**Changes:**
- Switched from `scratch` to `alpine:latest` base image (needed for netcat and wget)
- Installed `netcat-openbsd` for connection checking
- Installed `wget` for Docker healthcheck probes
- Installed `ca-certificates` for SSL/TLS support
- Added `wait-for-it.sh` script to container
- New entrypoint: `/app/wait-for-it.sh db:5432 -- /app/main`

### 3. Enhanced Docker Compose (`docker-compose.yml`)
**Core Service Changes:**
- Added healthcheck using `/v1/healthcheck` endpoint
- Healthcheck config:
  - Interval: 10s
  - Timeout: 5s
  - Retries: 5
  - Start period: 30s (grace period for startup)
- Added `restart: unless-stopped` for automatic recovery
- Added dependency on migration completion: `condition: service_completed_successfully`

**Migration Service Changes:**
- Added proper dependency on database health
- Added `restart: on-failure` for retry on errors
- Removed deprecated `links` directive

### 4. Go Service Improvements (`cmd/api/main.go`)
**RabbitMQ Connection with Retry Logic:**
- Added retry loop (10 attempts, 5s delay between attempts)
- Graceful degradation: Service starts even if RabbitMQ is unavailable
- Structured logging for connection attempts
- Continues accepting HTTP requests even without RabbitMQ

**Benefits:**
- Service doesn't crash if RabbitMQ isn't ready immediately
- Clear logging of connection status
- HTTP endpoints work independently of message queue

## Usage

### Starting Services
```bash
# Standard startup
make docker-up

# Or with docker-compose directly
docker-compose up -d
```

### Rebuilding After Changes
```bash
# Rebuild with cache
make docker-rebuild

# Force rebuild (no cache)
make docker-rebuild-force
```

### Monitoring Health
```bash
# Check service health status
docker ps

# View core service logs
make docker-logs-core

# Check healthcheck endpoint directly
curl http://localhost:4000/v1/healthcheck
```

### Troubleshooting

#### Service Not Starting
1. Check logs: `make docker-logs-core`
2. Verify database is healthy: `docker ps` (should show "healthy")
3. Check wait-for script output in logs
4. Ensure RabbitMQ is running on tranquara-network

#### Healthcheck Failing
1. Verify service is binding to 0.0.0.0:4000
2. Check if firewall is blocking port 4000
3. Increase `start_period` in docker-compose.yml if startup is slow

#### RabbitMQ Connection Issues
- Service will start without RabbitMQ (graceful degradation)
- Check logs for retry attempts
- Verify RabbitMQ is on tranquara-network: `docker network inspect tranquara-network`
- Ensure AI service docker-compose is running (contains RabbitMQ)

## Startup Sequence

The new startup order is:

1. **Database (PostgreSQL)** - starts first
2. **Database healthcheck** - waits until `pg_isready` succeeds
3. **Migrations** - runs after database is healthy
4. **Migration completion** - waits for migrations to finish
5. **Wait-for-it script** - checks database port availability
6. **Go Application** - starts main server
7. **RabbitMQ connection** - attempts connection with retries
8. **Healthcheck probes** - Docker verifies service is ready
9. **Service marked healthy** - ready to accept traffic

## Files Modified

- ✅ `wait-for-it.sh` (new) - Startup coordination script
- ✅ `Dockerfile.yml` - Updated base image and dependencies
- ✅ `docker-compose.yml` - Added healthchecks and proper startup order
- ✅ `cmd/api/main.go` - Added RabbitMQ retry logic
- ✅ `Makefile` - Added force rebuild command

## Testing

After implementing these changes:

1. **Cold Start Test**: 
   ```bash
   make docker-down
   make docker-up
   # Wait 30 seconds
   curl http://localhost:4000/v1/healthcheck
   ```

2. **Immediate Request Test**:
   ```bash
   make docker-down
   make docker-up
   # Send request immediately (should not hang)
   curl http://localhost:4000/v1/healthcheck
   ```

3. **Health Status Test**:
   ```bash
   docker ps
   # core_service should show "healthy" status after ~30s
   ```

## Expected Behavior

✅ **Before Fix**: Requests hung on first startup, required manual restart  
✅ **After Fix**: Service is ready immediately, requests never hang

The service now properly waits for dependencies and only accepts traffic when truly ready to serve requests.
