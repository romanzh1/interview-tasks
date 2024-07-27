BIN_DIR=${CURDIR}/bin
DOCKER_DIR=${CURDIR}/build/dev
DOCKER_YML=${DOCKER_DIR}/docker-compose.yml
ENV_NAME="route-256-ws-8"
MIGRATE_DIR=${CURDIR}/migrations

-include ./build/dev/.env

.PHONY: compose-up
compose-up:
	docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} up -d

.PHONY: compose-down
compose-down: ## terminate local env
	docker-compose -p ${ENV_NAME} -f -f ${DOCKER_YML} stop

.PHONY: compose-rm
compose-rm: ## remove local env
	docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} rm -fvs

.PHONY: compose-rs
compose-rs: ## remove previously and start new local env
	make compose-rm
	make compose-up

.PHONY: install-goose
install-goose:
	GOBIN=${BIN_DIR} go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: goose-migrate-all
goose-migrate-all: goose-migrate-1 goose-migrate-2

.PHONY: goose-migrate-down-all
goose-migrate-down-all: goose-migrate-1-down goose-migrate-2-down

.PHONY: goose-migrate-1
goose-migrate-1:
	${BIN_DIR}/goose -dir ${MIGRATE_DIR} postgres "\
		host=localhost\
		user=${DOCKER_POSTGRES_USER_1}\
		dbname=${DOCKER_POSTGRES_DB_1}\
		password=${DOCKER_POSTGRES_PASSWORD_1}\
		port=${DOCKER_POSTGRES_HOST_PORT_1}\
		sslmode=disable" up

.PHONY: goose-migrate-2
goose-migrate-2:
	${BIN_DIR}/goose -dir ${MIGRATE_DIR} postgres "\
		host=localhost\
		user=${DOCKER_POSTGRES_USER_2}\
		dbname=${DOCKER_POSTGRES_DB_2}\
		password=${DOCKER_POSTGRES_PASSWORD_2}\
		port=${DOCKER_POSTGRES_HOST_PORT_2}\
		sslmode=disable" up

.PHONY: goose-migrate-1-down
goose-migrate-1-down:
	${BIN_DIR}/goose -dir ${MIGRATE_DIR} postgres "\
		host=localhost\
		user=${DOCKER_POSTGRES_USER_1}\
		dbname=${DOCKER_POSTGRES_DB_1}\
		password=${DOCKER_POSTGRES_PASSWORD_1}\
		port=${DOCKER_POSTGRES_HOST_PORT_1}\
		sslmode=disable" down

.PHONY: goose-migrate-2-down
goose-migrate-2-down:
	${BIN_DIR}/goose -dir ${MIGRATE_DIR} postgres "\
		host=localhost\
		user=${DOCKER_POSTGRES_USER_2}\
		dbname=${DOCKER_POSTGRES_DB_2}\
		password=${DOCKER_POSTGRES_PASSWORD_2}\
		port=${DOCKER_POSTGRES_HOST_PORT_2}\
		sslmode=disable" down
