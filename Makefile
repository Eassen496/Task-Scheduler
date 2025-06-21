build:
	@-docker compose build

up: build
	@-docker compose up -d

down:
	@-docker compose down

logs:
	@-docker compose logs -f

clean:
	@-docker compose down --rmi all --volumes

.PHONY: build, up, down, logs, clean