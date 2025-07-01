# Linkaja Makefile Templates
Phony: help coverage test build-container build-instance build-all

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)


TARGET_MAX_CHAR_NUM=20
## Show Help make command
help:
	@echo ""
	@echo "Usage:"
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
        helpMessage = match(lastLine, /^## (.*)/); \
        if (helpMessage) { \
            helpCommand = substr($$1, 0, index($$1, ":")-1); \
                        helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
                        printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
                } \
        } \
    	{ lastLine = $$0 }' $(MAKEFILE_LIST)

## Running Unit-test Code https://tcashsquad.atlassian.net/wiki/spaces/DEVOPS/pages/1993801811/Golang+Makefiles+Standarization#Test
test:
	@echo "Step Unit Test"

## Running Code Coverage https://tcashsquad.atlassian.net/wiki/spaces/DEVOPS/pages/1993801811/Golang+Makefiles+Standarization#Coverage
coverage:
	@echo "Step code coverage"

## Running Code Dependency Check https://tcashsquad.atlassian.net/wiki/spaces/DEVOPS/pages/1993801811/Golang+Makefiles+Standarization#Check
check: # this step just example
	go mod tidy
	go mod download
	go mod verify

## Running code at local to debug propose
run: # this step just example
	go run main.go

## Build Code to Binary Artifact https://tcashsquad.atlassian.net/wiki/spaces/DEVOPS/pages/1993801811/Golang+Makefiles+Standarization#Build
build: # this step just example
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-s -w" -o $(CI_PROJECT_DIR)/$(ARTIFACT_DIR)/$(CI_PROJECT_NAME)