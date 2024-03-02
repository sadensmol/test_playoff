.PHONY:up
up:
	docker compose -f docker-compose.yaml up

.PHONY:down
down:
	docker compose -f docker-compose.yaml down --remove-orphans --volumes
