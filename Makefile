server:
	go run main.go

pullpostgres:
	docker pull postgres:15-alpine

pullredis:
	docker pull redis:latest

createnetwork:
	docker network create go-network

runpostgres:
	docker run --name postgres --network go-network -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=pass -d postgres:15-alpine

runredis:
	docker run --name redis --network go-network -d redis:latest

postgres:
	docker exec -it postgres psql

createdb:
	docker exec -it postgres createdb -U root go_chat

dropdb:
	docker exec -it postgres dropdb go_chat

migrationinit:
	migrate create -ext sql -dir db/migration -seq init

migrateup:
	migrate -path db/migration -database "postgresql://root:pass@localhost:5433/go_chat?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:pass@localhost:5433/go_chat?sslmode=disable" -verbose down

sqlc:
	sqlc generate

build:
	docker build -t go_chat_app:latest .

run:
	docker run --name app1 --network go-network -p 8080:8080 go_chat_app:latest

composeup:
	docker-compose up

composedown:
	docker-compose down

rm:
	docker rm -f $(docker ps -aq)

rmi:
	docker rmi $(docker images -q)

stop:
	docker stop $(docker ps -q)

removenetworks:
	docker network prune -f

removevolume:
	docker-compose down -v
