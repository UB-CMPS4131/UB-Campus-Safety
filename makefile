include .envrc

## run/api: run the cmd/api application

.PHONY: run/api
run/api:
	@echo 'Running Application...'
	@go run ./cmd/web 
##-port=4000 -env=development -limiter-burst=5 -limiter-rps=2 -limiter-enabled=false -db-dsn=${UB_DB_DSN}
## @go run ./cmd/api/ -port=4000 -env=production -db-dsn=${COMMENTS_DB_DSN}

## db/psql: connect to the database using psql (terminal)
.PHONY: db/psql
db/psql: 
	psql ${UB_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY db/migrations/new:
	@echo 'creating migration fles for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${UB_DB_DSN} up

## apply all down db migrations
.PHONY: db/migrations/down
db/migrations/down:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${UB_DB_DSN} down

## apply all down db migrations
.PHONY: db/migrations/force
db/migrations/force:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${UB_DB_DSN} force ${num}