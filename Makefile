.PHONY: deps network build run up down clean

APP_NAME=my_app
NETWORK=avito_network
DOCKER_COMPOSE_FILE=compose.yaml
PORT=8080
CONTAINER_NAME=avito_container

deps:
	@echo "Installing dependencies..."
	go mod tidy

network:
	@echo "Creating network..."
	@docker network ls | grep -q $(NETWORK) || docker network create $(NETWORK)

build:
	@echo "Creating docker image..."
	docker build -t $(APP_NAME) -f Dockerfile .

run:
	@echo "Running docker image..."
	docker run --network $(NETWORK) -d --name $(CONTAINER_NAME) -p $(PORT):$(PORT) $(APP_NAME)

up:
	@echo "Starting services with docker-compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

down:
	@echo "Stopping docker image..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

clean:
	@echo "Cleaning..."
	-docker-compose -f $(DOCKER_COMPOSE_FILE) down --rmi all
	-docker volume rm $$(docker volume ls -q --filter name=$(NETWORK)_*)
	-docker network rm $(NETWORK)
