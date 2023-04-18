DSN="postgresql://postgres:notsecret@localhost:5432/postgres" 

.PHONY: sqlc

run: sqlc
	POSTGRES_DSN="user=postgres password=notsecret host=localhost port=5432 dbname=postgres sslmode=disable" gin

sqlc:
	sqlc generate

# MIGRATIONS

up:
	@goose -dir ./migrations postgres $(DSN) up

down:
	@goose -dir ./migrations postgres $(DSN) down

status:
	@goose -dir ./migrations postgres $(DSN) status

gooseinit:
	@goose -dir ./migrations create init sql

schema:
	@psql $(DSN) -c "\d"
	@psql $(DSN) -c "\d authors"
