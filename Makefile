.PHONY: startdb

startdb:
	@docker run -d --rm \
			--volume "./resources/pg/init/init.sql:/docker-entrypoint-initdb.d/init.sql" \
			--volume "./resources/pg/data:/usr/pgdata" \
			--name=forecastDB \
			-p 5544:5432 \
			-e POSTGRES_DB=forecast \
			-e POSTGRES_USER=forecast \
			-e POSTGRES_PASSWORD=forecast postgres:16-alpine

stopdb:
	@docker stop forecastDB || exit 0

build:
	GOOS=linux GOARCH=amd64 CGO_ENABlED=0 go build -o bin/gocast .

run-test:
	make startdb
	go run main.go \
	./scripts/convert.sh \