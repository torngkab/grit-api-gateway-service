run:
	go run main.go

private:
	go env -w GOPRIVATE=github.com/torngkab/*

swag:
	swag init