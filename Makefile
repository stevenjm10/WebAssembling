.PHONY: docker-build docker-run docker-stop up down

IMAGE = webassembling:latest

# Build the Docker image (compiles Go -> WebAssembly inside Docker)
docker-build:
	docker build -t $(IMAGE) .

# Run the container on http://localhost:8080
docker-run: docker-build
	docker run -d --rm --name webassembling -p 8080:8080 $(IMAGE)

docker-stop:
	docker stop webassembling

# Same with Docker Compose
up:
	docker compose up -d --build

down:
	docker compose down
