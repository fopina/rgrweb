dev:
	@go run -tags=dev,nomin main.go --no-auth

dev-auth:
	@go run -tags=dev,nomin main.go --token a:12345

test:
	@go test -cover ./...

testv:
	@go test -v ./...

gen:
	@go generate ./...
