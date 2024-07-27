DOCKER_DIR=${CURDIR}/build/dev
DOCKER_YML=${DOCKER_DIR}/docker-compose.yml
ENV_NAME="route-256-ws-7"

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
