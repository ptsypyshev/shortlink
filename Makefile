APP=shortlink
.PHONY: help
help: Makefile ## Show this help
	@echo
	@echo "Choose a command run in "$(APP)":"
	@echo
	@fgrep -h "##" $(MAKEFILE_LIST) | sed -e 's/\(\:.*\#\#\)/\:\ /' | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: build
build: ## Build an application
	@echo "Building ${APP} ..."
	mkdir -p build
	go build -o build/${APP} cmd/main.go

run: ## Run an application
	@echo "Starting ${APP} ..."
	go run cmd/main.go

clean: ## Clean a garbage
	@echo "Cleaning"
	go clean
	rm -rf build

lint: ## Check a code by golangci-lint
	@echo "Linter checking..."
	golangci-lint run ./...


docker-build: ## Build docker image from Dockerfile
	@echo "Building ${APP} container image..."
	docker-compose build

docker-run:  ## Run docker containers
	@echo "Run docker containers for ${APP} ..."
	docker-compose up

docker-start:  ## Start docker containers
	@echo "Start docker containers for ${APP} ..."
	docker-compose start

docker-stop:  ## Stop docker containers
	@echo "Stop docker containers for ${APP} ..."
	docker-compose stop

docker-rm:  ## Rm docker containers
	@echo "Rm docker containers for ${APP} ..."
	docker-compose rm