DB_URL=postgresql://root:secret@localhost:5432/lamoda_db?sslmode=disable

postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14.4-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root lamoda_db

dropdb:
	docker exec -it postgres14 dropdb lamoda_db

migratecreate:
	migrate create -ext sql -dir migrations -seq init_db

migrateup:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path migration -database "$(DB_URL)" -verbose down

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migratecreate migrateup migratedown test