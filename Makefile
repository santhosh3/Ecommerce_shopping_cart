build:
	@go build -o bin/ecom cmd/main.go

run: build
	@cd bin && ./ecom


git:
	@git add .
	@git commit -m "$(m)"
	@git push origin HEAD:ecom

# make git m="Your commit message here"

docker-build:
	@docker compose up -d --build

docker-build-server:
	@docker compose up -d --build server

docker-dev-build:
	@docker compose -f docker-compose.dev.yaml up --build
