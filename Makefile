include .env
export

RED="\\033[31m"
GREEN="\\033[32m"
RESET="\\033[0m"

define delete-docker-service
	@if docker inspect ${APP_NAME}$1 >/dev/null 2>&1; then \
		docker stop ${APP_NAME}$1 >/dev/null 2>&1; \
		docker rm ${APP_NAME}$1 >/dev/null 2>&1; \
		echo -e "${GREEN}${APP_NAME}$1 deleted${RESET}"; \
	else \
		echo -e "${RED}${APP_NAME}$1 does not exist${RESET}"; \
	fi
endef

all: clear network-create db-create migrate-up build run
.PHONY: all

migrate-up: ### migration up
	@OUT=$$(goose -dir migrations postgres '${PG_URL}' up 2>&1); \
	if echo "$$OUT" | grep "successfully migrated" >/dev/null 2>&1; then \
	  	echo -ne "${GREEN}${APP_NAME}db "; \
	  	echo "$$OUT" | awk 'END {for (i=4; i<=NF; i++) printf "%s%s", $$i, (i<NF ? " " : "")}'; \
		echo -e "${RESET}"; \
	elif echo "$$OUT" | grep "no migrations to run" >/dev/null 2>&1; then \
		echo -ne "${GREEN}${APP_NAME}db "; \
	  	echo "$$OUT" | awk 'END {for (i=4; i<=NF; i++) printf "%s%s", $$i, (i<NF ? " " : "")}'; \
		echo -e "${RESET}"; \
	else \
		echo -e "${RED}${APP_NAME}db migration error"; \
		echo -e "$$OUT${RESET}"; \
	fi
.PHONY: migrate-up

migrate-down: ### migration down
	@OUT=$$(goose -dir migrations postgres '$(PG_URL)' down 2>&1); \
	if echo "$$OUT" | grep "OK" >/dev/null 2>&1; then \
	  	echo -ne "${GREEN}${APP_NAME}db migrated down to "; \
	  	echo "$$OUT" | awk 'END {for (i=4; i<=4; i++) printf "%s%s", $$i, (i<NF ? " " : "")}'; \
		echo -e "${RESET}"; \
	else \
		echo -e "${RED}${APP_NAME}db migration error"; \
		echo -e "$$OUT ${RESET}"; \
	fi
.PHONY: migrate-down

db-create: ### Creating postgres db docker instance
	@if docker inspect ${APP_NAME}db >/dev/null 2>&1; then \
  		echo -e "${RED}${APP_NAME}db already exists${RESET}"; \
  	else \
		docker run -p ${PG_PORT}:5432 --network=${APP_NETWORK}-network --restart unless-stopped --name ${APP_NAME}db -e POSTGRES_PASSWORD=${PG_PASSWORD} -v ${APP_NAME}_postgres_data:/var/lib/postgresql/data -d postgres:17-alpine >/dev/null 2>&1; \
		docker exec ${APP_NAME}db sh -c 'until pg_isready -U postgres; do sleep 0.5; done' >/dev/null 2>&1; \
		docker exec ${APP_NAME}db psql -U postgres -tAc "SELECT 1 FROM pg_database WHERE datname='${APP_NAME}db'" | grep -q 1 || docker exec ${APP_NAME}db psql -U postgres -c "CREATE DATABASE ${APP_NAME}db" >/dev/null 2>&1; \
		echo -e "${GREEN}${APP_NAME}db created${RESET}"; \
	fi
.PHONY: db-create

db-delete: ### Deleting postgres db docker instance
	$(call delete-docker-service,db)
.PHONY: db-delete

build: ### Building app
	@echo -e "${GREEN}${APP_NAME} is building...${RESET}"; \
	OUT=$$(go build -o ${APP_NAME} ./cmd/app 2>&1); \
	EXIT_CODE=$$?; \
	if [ $$EXIT_CODE -eq 0 ]; then \
	  	echo -e "${GREEN}${APP_NAME} builded successfully${RESET}"; \
	else \
	  	echo -ne "${RED}${APP_NAME} building error: "; \
	  	echo -e "$$OUT${RESET}"; \
	fi
.PHONY: build

run: ### Running app
	@./${APP_NAME}
.PHONY: run

network-create:
	@if docker network ls | grep "${APP_NETWORK}-network" >/dev/null 2>&1; then \
  		echo -e "${RED}${APP_NETWORK}-network already exists${RESET}"; \
  	else \
  	  	docker network create ${APP_NETWORK}-network >/dev/null 2>&1; \
  	  	echo -e "${GREEN}${APP_NETWORK}-network created${RESET}"; \
  	fi
.PHONY: network-create

network-delete:
	@if docker network inspect ${APP_NETWORK}-network >/dev/null 2>&1; then \
  		docker network rm "${APP_NETWORK}-network" >/dev/null 2>&1; \
  		echo -e "${GREEN}${APP_NETWORK}-network deleted${RESET}"; \
  	else \
  	  	echo -e "${RED}${APP_NETWORK}-network does not exist${RESET}"; \
  	fi
.PHONY: network-delete

clear: db-delete network-delete ### Cleaning up
.PHONY: clear
