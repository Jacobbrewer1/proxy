# Define variables
hash = $(shell git rev-parse --short HEAD)
registry = ghcr.io/jacobbrewer1/proxier
DATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

summary:
	@echo "Building proxier"
	make linux
	chmod +x bin/app
PR-approval:
	@echo "Running PR CI"
	go build ./...
	go vet ./...
	go test ./...
ci: docker
	# Push the image to the registry
	docker push $(registry):$(hash)
linux: clean
	@echo "Building for linux"
	GOOS=linux GOARCH=amd64 go build -o bin/app -ldflags '-X main.Commit=$(hash) -X main.Date=$(DATE)'
docker: linux
	@echo "Building docker image"
	# Build the docker image
	docker build -t $(registry):$(hash) .
clean:
	@echo "Cleaning up"
	# Remove the bin directory
	rm -rf bin
