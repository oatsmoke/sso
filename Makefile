include .env
export
gen:
	protoc \
		-I api api/sso/v1/*.proto \
		--go_out=./internal/gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=./internal/gen \
		--go-grpc_opt=paths=source_relative

migrate-up:
	goose -dir ${MIGRATION_PATH} ${DRIVER} ${DB_STRING} up

migrate-down:
	goose -dir ${MIGRATION_PATH} ${DRIVER} ${DB_STRING} down