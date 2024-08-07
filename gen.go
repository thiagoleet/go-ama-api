package gen

//go:generate go run ./cmd/tools/terndotenv/main.go
//go:geneatate sqlc generate -f ./internal/store/pgstore/sqlc.yml
