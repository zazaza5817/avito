.PHONY: docker-up docker-down lint deps

docker-up: ## Запустить все через Docker Compose
	docker-compose up --build -d

docker-down: ## Остановить все
	docker-compose down

lint: ## Запустить линтер
	golangci-lint run

fmt: ## Форматировать код
	go fmt ./...
	goimports -w .

deps: ## Установить зависимости
	go mod download
	go mod tidy

load-test: ## Запустить нагрузочное тестирование с помощью k6
	k6 run k6-load-test.js
