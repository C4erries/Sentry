APP_NAME := sentry
CMD_DIR := ./cmd/$(APP_NAME)
BINARY := $(APP_NAME)

DOCKER_COMPOSE := docker-compose
GOCMD := go

.PHONY: build
build:
	$(GOCMD) build -o $(BINARY) $(CMD_DIR)

.PHONY: run
run:
	./$(BINARY)

.PHONY: clean
	rm -f $(BINARY)

.PHONY: docker-build
docker-build:
	docker build -t $(APP_NAME):latest .

.PHONY: up
up:
	$(DOCKER_COMPOSE) up --build "sentry"

.PHONY: produce
produce:
	$(DOCKER_COMPOSE) exec sentry /sentry produce --type=login --user_id=123 --count=10 --country=RU --ip=0.0.0.52 --method=post --success=true
.PHONY: migrate
	$(DOCKER_COMPOSE) exec sentry /sentry migrate