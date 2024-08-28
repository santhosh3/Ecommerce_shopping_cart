build:
	@go build -o bin/ecom cmd/main.go

run: build
	@./bin/ecom

git:
	@git add .
	@git commit -m "commit_changes"
	@git push

docker_build:
	@docker build -t ecom .