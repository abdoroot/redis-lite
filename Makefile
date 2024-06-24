run:build
	@./temp/lb
build:
	@go build -o temp/lb ./cmd/main.go
set:
	@go run scripts/client_main.go set name "mohamed" EXP 60
get:
	@go run scripts/client_main.go get name	