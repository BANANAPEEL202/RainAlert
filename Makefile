all: build

validate:
	@go run ./cmd/validate-config/main.go

build: validate
	GOOS=linux GOARCH=amd64 go build -o cmd/lambda/bootstrap cmd/lambda/main.go
	@echo "✅ Build complete."
	@echo "======================================="
	@echo ""
	@echo "Next steps to deploy:"
	@echo "  - Guided deploy:   sam deploy --guided"
	@echo "  - Standard deploy: sam deploy"
	@echo ""