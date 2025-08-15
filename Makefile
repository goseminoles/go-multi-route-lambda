# ================================
# Go Lambda Makefile for provided.al2 (ARM64)
# ================================

APP_NAME := bootstrap
ZIP_NAME := $(APP_NAME).zip

GOOS := linux
GOARCH := arm64        # <--- Use ARM64 architecture
LAMBDA_FUNCTION := go-multi-route-lambda
REGION := us-east-1

.PHONY: all clean build zip local deploy

all: clean build zip

clean:
	rm -f $(APP_NAME) $(ZIP_NAME)

build:
	@echo "Building $(APP_NAME) for Lambda provided.al2 runtime (ARM64)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(APP_NAME) main.go

zip: build
	@echo "Creating zip $(ZIP_NAME)..."
	zip -r $(ZIP_NAME) $(APP_NAME)

local:
	@echo "Running locally at http://localhost:8080..."
	go run main.go


deploy: zip
	@echo "Deploying $(ZIP_NAME) to Lambda function $(LAMBDA_FUNCTION)..."
	aws lambda update-function-code \
		--function-name $(LAMBDA_FUNCTION) \
		--zip-file fileb://$(ZIP_NAME) \
		--region $(REGION)
	@echo "Deployment complete. Invoke via API Gateway."
