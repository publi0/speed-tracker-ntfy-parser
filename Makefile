# Set the base go command
GO=go

# Set the Golang binary output name
BINARY_NAME=ntfy-parser

# Set the Golang binary output directory
BINARY_UNIX=$(BINARY_NAME)

# Set the default goal to be "all" if no goal is provided
.DEFAULT_GOAL := all

# Build the Golang application
build:
	$(GO) build -o $(BINARY_NAME) -v

# Test the Golang application
test:
	$(GO) test -v ./...

# Clean up the Golang application
clean:
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run the Golang application
run:
	$(GO) run .

# Build the Golang application for unix
build-unix:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY_UNIX) -v

# Default make goal
all: test build

docker-deploy: build-unix
	docker build -t ntfy-parser-speed-tracker . &&
	docker tag ntfy-parser-speed-tracker:latest felipecanton/ntfy-parser-speed-tracker:latest &&
	docker push felipecanton/ntfy-parser-speed-tracker:latest
