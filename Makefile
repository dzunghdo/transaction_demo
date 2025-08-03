SRC_PATH:= ${PWD}

mod:
	@go mod tidy
	@go mod vendor

build:
	@go build -o srv ${SRC_PATH}/cmd/srv/...

dev-tools:
	@go get -u -v github.com/swaggo/swag/cmd/swag@v1.16.3
	@go get -u -v github.com/golang/mock/gomock@v1.6.0
	@go get -u -v golang.org/x/tools/cmd/goimports
	@go install github.com/pressly/goose/v3/cmd/goose@v3
	@export GOOSE_MIGRATION_DIR='db/migrations'
	@export GOOSE_DRIVER=postgres

migrate-up:
	@goose postgres -dir db/migrations "postgresql://postgres:root123@localhost:15432/example_db?sslmode=disable" up

migrate-down:
	@goose postgres -dir db/migrations "postgresql://postgres:root123@localhost:15432/example_db?sslmode=disable" down

swag:
	@echo '$(shell swag --version)'
	@swag init -g app/interface/api/route/route.go --parseVendor true --exclude db,deployment,scripts,vendor

mock:
	@go generate ./...