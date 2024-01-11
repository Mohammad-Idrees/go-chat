server:
	go run main.go

pullpostgres:
	docker pull postgres:15-alpine

postgresinit:
	docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=pass -d postgres:15-alpine

postgres:
	docker exec -it postgres15 psql

createdb:
	docker exec -it postgres15 createdb -U root go_chat

dropdb:
	docker exec -it postgres15 dropdb go_chat

migrationinit:
	migrate create -ext sql -dir db/migration -seq init

migrateup:
	migrate -path db/migration -database "postgresql://root:pass@localhost:5433/go_chat?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:pass@localhost:5433/go_chat?sslmode=disable" -verbose down

sqlc:
	sqlc generate


