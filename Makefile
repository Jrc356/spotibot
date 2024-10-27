APP_NAME=spotibot
VERSION=v0.0.1

OUT=out
go-build-dev:
	go build -gcflags=all="-N -l" -o "$(OUT)/$(APP_NAME)" ./cmd/**

go-build-prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

go-run:
	go run ./cmd/main.go

docker-build-prod:
	docker build -t "$(APP_NAME):$(VERSION)" .

DEV_TAG="dev"
docker-build-dev:
	docker build --target dev -t "$(APP_NAME):dev" .

DELV_PORT=2345
docker-run-dev: docker-build-dev
	docker run --env-file .env -p $(DELV_PORT):$(DELV_PORT) -v ./:/app $(APP_NAME):dev

run-dev: docker-build-dev docker-run-dev
