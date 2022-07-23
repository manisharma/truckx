up: ## bring app and all its supporting services up
	docker-compose up --build

down: ## shutdown app and all its supporting services
	docker-compose down

create-network: ## create a separate network bridge for the app
	docker network create truckx_network

create-postgres: ## create a postgres databse
	docker run --name truckx_postgres -d -p 5432:5432 -e "POSTGRES_USER=admin" -e "POSTGRES_PASSWORD=passw0rd" -e "POSTGRES_DB=postgres" --net truckx_network postgres:latest

create-postgres-explorer: ## create a postgres-explorer databse in the network
	docker run -p 5050:80 -e "PGADMIN_DEFAULT_EMAIL=admin@truckx.com" -e "PGADMIN_DEFAULT_PASSWORD=passw0rd" --net truckx_network-d dpage/pgadmin4

migrate: ## run db migration
	export PGURL="postgres://admin:passw0rd@localhost:5432/postgres?sslmode=disable"
	migrate -database $(PGURL) -path internal/db/migrations/ down
	migrate -database $(PGURL) -path internal/db/migrations/ up

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help