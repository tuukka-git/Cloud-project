#!/bin/bash

# Define PostgreSQL Docker image version and container details
POSTGRES_VERSION="16"
CONTAINER_NAME="my_postgres"
POSTGRES_USER="myuser"
POSTGRES_PASSWORD="mypassword"
POSTGRES_DB="mydb"
POSTGRES_PORT="5432"

# Function to check if the Docker image exists
function check_postgres_image() {
    docker images | grep -q "postgres.*$POSTGRES_VERSION"
    return $?
}

# Function to pull the Docker image if needed
function pull_postgres_image() {
    echo "Pulling PostgreSQL Docker image (postgres:$POSTGRES_VERSION)..."
    docker pull postgres:$POSTGRES_VERSION
}

# Function to check if the container is already running
function check_container_running() {
    docker ps | grep -q "$CONTAINER_NAME"
    return $?
}

# Function to start the PostgreSQL container
function start_postgres_container() {
    echo "Starting PostgreSQL container ($CONTAINER_NAME)..."
    docker run --name "$CONTAINER_NAME" -e POSTGRES_USER="$POSTGRES_USER" -e POSTGRES_PASSWORD="$POSTGRES_PASSWORD" -e POSTGRES_DB="$POSTGRES_DB" -p "$POSTGRES_PORT":5432 -d postgres:$POSTGRES_VERSION
}

# Check if Docker is installed
if ! [ -x "$(command -v docker)" ]; then
    echo "Error: Docker is not installed." >&2
    exit 1
fi

# Check if the PostgreSQL image is already available
if check_postgres_image; then
    echo "PostgreSQL Docker image (postgres:$POSTGRES_VERSION) is already present."
else
    echo "PostgreSQL Docker image (postgres:$POSTGRES_VERSION) not found."
    pull_postgres_image
fi

# Check if the container is already running
if check_container_running; then
    echo "PostgreSQL container ($CONTAINER_NAME) is already running."
else
    echo "PostgreSQL container ($CONTAINER_NAME) is not running. Starting it now..."
    start_postgres_container
fi

# Export PostgreSQL connection string

#export DATABASE_URL="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:$POSTGRES_PORT/$POSTGRES_DB"
echo DATABASE_URL="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:$POSTGRES_PORT/$POSTGRES_DB"

go run main.go