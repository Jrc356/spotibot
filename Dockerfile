FROM golang:1.23.2 AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# mount the repo as a volume before running this target
# -v ./:/app
FROM base AS dev
RUN go install github.com/air-verse/air@v1.61.1
ENTRYPOINT ["air", "-c", ".air.toml"]

FROM base AS build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

FROM scratch AS production
WORKDIR /app
COPY --from=build /app/main .
ENTRYPOINT ["/app/main"]
