build:
	@go build -o bin/ecom cmd/main.go

run: build
	@./bin/ecom


git:
	@git add -- :!Dockerfile
	@git commit -m "$(m)"
	@git push

# make git m="Your commit message here"

docker-build:
	@docker compose up -d --build

docker-build-server:
	@docker compose up -d --build server

docker-dev-build:
	@docker compose -f docker-compose.dev.yaml up --build
