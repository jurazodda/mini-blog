# Запуск приложения с нужными переменными окружения
run:
	@mkdir -p logs
	@export DSN_URL="host=localhost user=postgres password=postgres dbname=mini-blog port=6543 sslmode=disable"; \
	export SIGNING_KEY="jfsgvkje9o9309-30"; \
	go run main.go

# Запуск всей инфраструктуры через Docker Compose
up:
	docker compose -f ./docker-compose-local.yml up -d

# Остановка всей инфраструктуры
down:
	docker compose -f ./docker-compose-local.yml down

# Запуск линтера (golangci-lint)
lint:
	golangci-lint run

# Запуск тестов с генерацией отчёта о покрытии
test:
	go test -cover ./...
