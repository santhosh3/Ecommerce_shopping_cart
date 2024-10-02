build:
	@go build -o bin/ecom cmd/main.go

stopProcess:
	@echo "Stopping process on port 50051..."
	@sudo fuser -k 50051/tcp || true

run: stopProcess build
	@cd bin && ./ecom

git:
	@git add .
	@git commit -m "$(m)"
	@git push origin HEAD:ecom

docker-build:
	@docker compose up -d --build

docker-build-server:
	@docker compose up -d --build server

docker-dev-build:
	@docker compose -f docker-compose.dev.yaml up --build

grpc:
	@protoc --go_out=. --go-grpc_out=. proto/service.proto
