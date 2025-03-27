# Transactional Outbox Pattern Implementation

## How it works

### File Structure
```
.
├── cmd
│   └── example_client        # Example client application
├── src
│   ├── config                # Configuration for DB and NATS
│   ├── internal
│   │   ├── client            # Client connections (NATS, Postgres)
│   │   ├── model             # Data models
│   │   ├── obj               # Message objects
│   │   ├── repository        # Repository implementations
│   │   │   └── mocks         # Mock repositories for testing
│   │   ├── service           # Core services implementation
│   │   └── util              # Utility functions
│   └── outbox                # Public API for the library
├── Dockerfile                # Docker configuration
├── docker-compose.yml        # Docker compose configuration
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── Makefile                  # Build and run commands
```

- This is a library implementation of Transactional Outbox pattern
- There are two major Service components in this Structure
    - PostMessage Service: allows clients to post messages, which are persisted in the Postgres DB
    - Outbox Message Service: allows clients to start a separate BatchJob, that continuously Pushes unread messages to NATS server.
- Both of these major Services have been exposed via external Constructors that can be accesed by clients
- Clients need to provide configuration details to connect with Postgres DB and NATS server, and then they can easily call these lib services to serve the OutboxPattern
- An example client has been created in `cmd/example_client/main.go` to run the application and to show how the Library works.

## Running the Application

This application uses Docker and can be managed using the provided Makefile commands.

### Available Commands

- `make build`: Builds the Docker containers for the application
- `make run`: Starts the application in detached mode
- `make down`: Stops all running containers
- `make clean`: Stops containers and removes unused Docker resources
- `make logs`: Displays container logs in follow mode
- `make re-run`: Stops containers, cleans resources, and restarts the application
- `make test`: Runs the test suite
- `make test-coverage`: Generates a test coverage report
- `make all`: Builds and runs the application (default command)

### Getting Started

1. To build and start the application for the first time:
   ```
   make all
   ```

2. To view application logs:
   ```
   make logs
   ```

3. To stop the application:
   ```
   make down
   ```

### Development Workflow

1. Make changes to the codebase
2. Rebuild and restart the application:
   ```
   make re-run
   ```

3. Run tests to verify your changes:
   ```
   make test
   ```

4. Check test coverage:
   ```
   make test-coverage
   ```

5. Create mocks using mockery
   ```
   mockery --name=IOutboxMessageRepository \
   --dir=./src/internal/repository \
   --output=./src/internal/repository/mocks \
   --filename=outbox_message_repository_mock.go
   ```

