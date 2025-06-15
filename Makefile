run:
	@go run main.go
up:
	docker compose -f ./docker-compose-local.yml up -d
down:
	docker compose -f ./docker-compose-local.yml down
