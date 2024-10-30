.PHONY: startdb

startdb:
	@docker run --network gocast-development_google -d --rm \
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

test: # Use "go tool cover -html=./bin/metrictests.out" in terminal to open coverage in web (Open it in Chrome, Firefox doesn't work because it is not based on chromium)
	go test ./pkg/metric/ ./pkg/utils/ -cover -coverprofile ./bin/metrictests.out

run:
	docker-compose up

run2:
	go run main.go

stop:
	docker-compose down