# tests
test-databases:
	go test ./databases/redisIO.go ./databases/redisIO_test.go
# development
dev: export SERVER_PORT=5000
dev:
	go run main.go
# production
prod: export GIN_MODE=release
prod: export SERVER_PORT=5000
prod:
	go run main.go
