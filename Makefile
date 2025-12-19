.PHONY: cli
cli: ## Builds CLI Command
# 	GOOS=darwin COARCH=amd64 go build -o cli ./cmd/cli
	go build -o cli ./cmd/cli

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.DEFAULT_GOAL := help