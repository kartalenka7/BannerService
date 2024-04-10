banner-service-up: ## Create and run app containers
	docker-compose up --build banner-service

banner-service-down: ## Stop and remove app containers
	docker compose --file docker-compose.yml down -v