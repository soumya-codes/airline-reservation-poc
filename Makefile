.PHONY: test, generate-init-script

# Define Docker Compose command
DOCKER_COMPOSE := docker-compose

# Docker Compose file
COMPOSE_FILE := ./deployment/docker-compose.yml

#Generate SQL
generate-sql:
	sqlc generate

sleep:
	@echo "Sleeping for 3 seconds..."
	sleep 3
	@echo "Waking up..."

# Run Docker containers
setup: teardown sleep
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up -d

# Remove Docker containers and associated volumes
teardown:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) down -v

# View logs of Docker containers
logs:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) logs -f

# Run tests
test: setup
	go test -v ./...

