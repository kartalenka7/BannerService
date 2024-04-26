banner-service-up: ## Create and run app containers
	docker-compose up --build banner-service

banner-service-down: ## Stop and remove app containers
	docker compose --file docker-compose.yml down -v

## Test:
test: ## Run tests
	@docker-compose --file docker-compose-test.yml up -d --force-recreate
	@go test -count=1 -v ./tests
	@go test -bench=. -benchtime=1s ./tests
	@docker-compose --file docker-compose-test.yml down -v