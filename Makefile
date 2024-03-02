.PHONY:run
run: down
	docker compose -f docker-compose.yaml --profile run up

.PHONY:up
up:
	docker compose -f docker-compose.yaml up

.PHONY:down
down:
	docker compose -f docker-compose.yaml down --remove-orphans --volumes

.PHONY:test
test:
	go test -v ./...

