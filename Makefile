postgres:
	sudo docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine

createbd: 
	sudo docker exec -it postgres createdb --username=root --owner=root course
dropdb:
	sudo docker exec -it postgres dropdb  course

migrateup: 
	migrate -path migration/ -database "postgresql://root:secret@localhost:5432/course?sslmode=disable" -verbose up

migratedown:
	migrate -path migration/ -database "postgresql://root:secret@localhost:5432/course?sslmode=disable" -verbose down

createredis:
	docker run --name test2 -p 6380:6379 -d redis

proto:
	protoc --go_out=./pkg/course --go_opt=paths=source_relative \
    --go-grpc_out=./api/courses --go-grpc_opt=paths=source_relative \
    api/courses/service.proto



.PHONY: postgres createbd dropdb migrateup migratedown createredis
