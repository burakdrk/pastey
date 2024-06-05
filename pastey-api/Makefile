postgres:
	docker run --name pastey-db --network pastey-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it pastey-db createdb --username=root --owner=root pastey

dropdb:
	docker exec -it pastey-db dropdb pastey

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/pastey?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/pastey?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/pastey?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/pastey?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createdb dropdb postgres migrateup migratedown sqlc server api migrateup1 migratedown1 test