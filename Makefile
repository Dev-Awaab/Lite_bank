postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root lite_bank

dropdb:
	docker exec -it postgres12 dropdb lite_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/lite_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/lite_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

startdb:
	docker start postgres12

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc startdb test
