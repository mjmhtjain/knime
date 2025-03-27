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