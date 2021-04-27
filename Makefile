GO_SRCS := $(shell find . -type f -name '*.go')

count-go: ## Count number of lines of all go codes.
	find . -name "*.go" -type f | xargs wc -l | tail -n 1

docker-build: $(GO_SRCS)
	bash ./scripts/build.sh

go-build: $(GO_SRCS)
	go build ./

go-test: $(GO_SRCS)
	go test -v
