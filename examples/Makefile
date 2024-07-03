default: down build up

bin/%: %/main.go
	GOOS=linux CGO_ENABLED=0 go build -o "$@" "$<"

build: \
	bin/proxy \
	bin/client \
	bin/server \
	bin/client-as-server \
	bin/server-as-client

up:
	docker compose up

down:
	docker compose down
