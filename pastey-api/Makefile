postgres:
	docker run --name pastey-db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it pastey-db createdb --username=root --owner=root pastey

dropdb:
	docker exec -it pastey-db dropdb pastey

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/pastey?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/pastey?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: createdb dropdb postgres migrateup migratedown sqlc