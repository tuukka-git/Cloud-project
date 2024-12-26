# Makefile

# Variables
DOCKER_COMPOSE_FILE = docker-compose.yml
REACT_PROJECT_DIR = ./frontend/team-gen-frontend
REACT_START_CMD = npm start # Change to 'yarn start' if you use Yarn
WAIT_TIME = 10  # Number of seconds to wait before starting the React project

.PHONY: up react start clean down

# Start Docker Compose services
up:
	@echo "Starting Docker Compose services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "Waiting for Docker Compose services to be ready..."
	sleep $(WAIT_TIME)

# Start the React project
react:
	@echo "Starting React project..."
	cd $(REACT_PROJECT_DIR) && $(REACT_START_CMD)

# Start both Docker Compose and React project sequentially
start: up react

# Stop and remove Docker Compose services
down:
	@echo "Stopping Docker Compose services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Clean up stopped containers and images
clean:
	@echo "Cleaning up Docker resources..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down --rmi all --volumes --remove-orphans
