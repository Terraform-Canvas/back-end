.PHONY: clean critic security lint test build run
include .env
APP_NAME = apiserver
BUILD_DIR = $(PWD)/build
MIGRATIONS_FOLDER = $(PWD)/platform/migrations
DATABASE_URL = mysql://root:${DB_PASSWORD}@tcp(localhost:3306)/mysql?query

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./...

security:
	gosec ./...

lint:
	golangci-lint run ./...

test: clean critic security lint
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

build: test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: swag build
	$(BUILD_DIR)/$(APP_NAME)

migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" up

migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" down

migrate.force:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" force $(version)

docker.run: docker.network docker.mysql swag docker.fiber migrate.up

docker.network:
	docker network inspect dev-network >/dev/null 2>&1 || \
	docker network create -d bridge dev-network

docker.fiber.build:
	docker build -t fiber .

docker.fiber: docker.fiber.build
	docker run --rm -d \
		--name cgapp-fiber \
		--network dev-network \
		-p 3000:3000 \
		fiber

docker.mysql:
	docker run --rm -d \
		--name cgapp-mysql \
		--network dev-network \
		-e MYSQL_ROOT_PASSWORD=${DB_PASSWORD} \
		-e MYSQL_DATABASE=mysql \
		-v ${HOME}/dev-mysql:/var/lib/mysql \
		-p 3306:3306 \
		mysql

docker.postgres:
	docker run --rm -d \
		--name cgapp-postgres \
		--network dev-network \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=postgres \
		-v ${HOME}/dev-postgres/data/:/var/lib/postgresql/data \
		-p 5432:5432 \
		postgres

docker.redis:
	docker run --rm -d \
		--name cgapp-redis \
		--network dev-network \
		-p 6379:6379 \
		redis

docker.stop: docker.stop.fiber docker.stop.mysql #docker.stop.postgres docker.stop.redis

docker.stop.fiber:
	docker stop cgapp-fiber

docker.stop.postgres:
	docker stop cgapp-postgres

docker.stop.redis:
	docker stop cgapp-redis

docker.stop.mysql:
	docker stop cgapp-mysql

swag:
	swag init
