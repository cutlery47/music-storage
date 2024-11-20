run:
	go run cmd/main.go

up:
	docker compose up -d

build:
	docker compose build

up_build:
	docker compose up -d --build

unbuild:
	docker rmi music-postgres-image
	docker rmi music-app-image

stop:
	docker compose stop

down:
	docker compose down

clean:
	docker volume rm music-storage_postgres_data