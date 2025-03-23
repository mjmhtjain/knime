.PHONY: build run clean re-run all logs test test-coverage

# Docker commands
DOCKER=docker
DOCKER_COMPOSE=docker-compose

# Build the application using Docker
build:
	$(DOCKER_COMPOSE) build

# Run the application with Docker
run:
	$(DOCKER_COMPOSE) up -d

# Clean Docker artifacts
clean:
	$(DOCKER_COMPOSE) down
	$(DOCKER) system prune -f

# Show logs in follow mode
logs:
	$(DOCKER_COMPOSE) logs -f

# Re-run the application
re-run: clean run

# Run unit tests
test:
	$(DOCKER_COMPOSE) run --rm app test

# Generate test coverage report
test-coverage:
	$(DOCKER_COMPOSE) run --rm app test --coverage

# Default target
all: build run 