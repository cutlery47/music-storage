up:
	docker compose up -d

build:
	docker compose build

unbuild:
	docker rmi music-postgres-image
	docker rmi music-app-image

stop:
	docker compose stop

down:
	docker compose down