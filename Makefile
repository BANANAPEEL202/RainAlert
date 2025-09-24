validate:
	go run ./cmd/validate-config/main.go

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cmd/lambda/bootstrap cmd/lambda/main.go
