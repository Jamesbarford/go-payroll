ENTRY := ./main.go
TARGET := ./api.out
CC := go
MIGRATION_PATH := ./database/migrations

# Modifiyable variables
DB_USER?=""
DB_PASSWORD?=""
DB_NAME?="$(whoami)"
DB_HOST?="localhost"
DB_PORT?=5432

.PHONY: set_envvars migrate_up clean

all: build

build:
	$(CC) build -o $(TARGET) $(ENTRY)

clean:
	rm $(TARGET)

test:
	go test ./...

migrate_up:
	migrate -path $(MIGRATION_PATH) \
		-database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
		-verbose up

migrate_down:
	migrate -path $(MIGRATION_PATH) \
		-database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
		-verbose down

database_seed:
	psql -U $(DB_USER) -d $(DB_NAME) -a -f ./database/seed.sql
