APP_NAME="spotibot"
VERSION="0.0.1"

go-build:
	go build -o out/$(APP_NAME) ./cmd/**

go-run:
	go run ./cmd/main.go

docker-build-prod:
	docker build -t $(APP_NAME):$(VERSION) .

DEV_TAG="dev"
docker-build-dev:
	docker build --target dev -t $(APP_NAME):dev .

docker-run-dev: docker-build-dev
	docker run -v ./:/app $(APP_NAME):dev

run-dev: docker-build-dev docker-run-dev
