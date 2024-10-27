FROM golang:1.23.2 AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev
RUN go install github.com/air-verse/air@v1.61.1 && \
  go install github.com/go-delve/delve/cmd/dlv@v1.23.1
# delve
EXPOSE 2345
ENTRYPOINT ["air", "-c", ".air.toml"]

FROM base AS build
COPY . .
RUN make go-build-prod

FROM scratch AS production
WORKDIR /app
COPY --from=build /app/main .
ENTRYPOINT ["/app/main"]
