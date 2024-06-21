.PHONY: test generate-sql setup optional-teardown teardown logs sleep

# Define Docker Compose command
DOCKER_COMPOSE := docker-compose

# Docker Compose file
COMPOSE_FILE := ./deployment/docker-compose.yml

# Optional variables
TEST_OUTPUT_FILE ?=
RUN_TEARDOWN ?= 1

# Generate SQL
generate-sql:
	sqlc generate

# Sleep for 3 seconds
sleep:
	@echo "Sleeping for 3 seconds..."
	sleep 3
	@echo "Waking up..."

# Run Docker containers
setup: teardown sleep
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up -d

# Optional teardown step
optional-teardown:
ifneq ($(RUN_TEARDOWN), 0)
	$(MAKE) teardown
else
	@true
endif

# Remove Docker containers and associated volumes
teardown:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) down -v

# View logs of Docker containers
logs:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) logs -f

# Run tests
# Usage: make test TEST_OUTPUT_FILE=test_output.txt RUN_TEARDOWN=0
test: setup
ifneq ($(TEST_OUTPUT_FILE),)
	go test -v ./... > $(TEST_OUTPUT_FILE) 2>&1 || (result=$$?; $(MAKE) optional-teardown; exit $$result)
else
	go test -v ./... || (result=$$?; $(MAKE) optional-teardown; exit $$result)
endif
	$(MAKE) optional-teardown