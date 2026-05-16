up:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f

migrate-up:
	docker compose run --rm migrate -path /migrations -database 'postgres://postgres:postgres@postgres:5432/autopark?sslmode=disable' up

migrate-down:
	docker compose run --rm migrate -path /migrations -database 'postgres://postgres:postgres@postgres:5432/autopark?sslmode=disable' down 1
