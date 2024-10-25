.PHONY: build test docker run clean help lint

APP_SERVICE=image-previewer

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# Сборка приложения
build:
	go build -o bin/$(APP_SERVICE) ./cmd/$(APP_SERVICE)

# Запуск юнит-тестов
test:
	go test -race -count=100 ./...

# Сборка Docker-образа
docker:
	docker build -t $(APP_SERVICE):latest .

# Запуск Docker-контейнера
docker.run:
	docker run -d -p 8080:8080 --name $(APP_SERVICE) \
		-v $(PWD)/previews:/app/previews \
		$(APP_SERVICE):latest

# Остановка и удаление Docker-контейнера
docker.clean:
	docker rm -f $(APP_SERVICE) || true
	docker rmi $(APP_SERVICE):latest || true

lint:
	golangci-lint run -v ./...

lint-fix:
	gofumpt -l -w .
	golangci-lint run ./... --fix


run:
	docker-compose up --build $(APP_SERVICE)

integration-test:
	docker-compose up --build

down:
	docker-compose down