.PHONY: migrate-up migrate-down swagger
swagger:
	swag init -g cmd/main.go


migrate-up:
	docker compose exec db psql -U postgres -d tenderdb -f /docker-entrypoint-initdb.d/tables-up.sql

migrate-down:
	docker compose exec db psql -U postgres -d tenderdb -f /docker-entrypoint-initdb.d/tables-down.sql

run_db:
	docker compose build db && docker compose up -d db
	make migrate-up

run:
	docker compose build app && docker compose up -d app