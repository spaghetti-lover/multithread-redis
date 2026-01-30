# Test the application
test:
	go mod tidy
	go test ./... -v

# Test with vendor (slower but reproducible)
test-vendor:
	go mod tidy
	go mod vendor
	go test ./... -v -mod=vendor

# Test for CI
test-ci:
	go mod download
	go mod verify
	go test ./... -v

# Development
dev:
	docker-compose up --build

# Deployment (single container)
deploy:
	go mod tidy
	go mod vendor
	docker build -f Dockerfile.render -t myredis-app .
	docker run -p 3000:3000 -p 8080:8080 -p 6379:6379 myredis-app

# Document
doc:
	pkgsite -http :8080
# Clean up
clean:
	docker-compose down -v
	docker system prune -f