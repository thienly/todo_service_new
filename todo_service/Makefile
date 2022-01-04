all:
	@echo server generate database psql-start psql-stop migration-up migration-down
server:
	go run cmd/server/*.go
generate:
	@bash protoc-gen.sh
psql-start:
	docker run --name todo-postgres -p 5432:5432 -e POSTGRES_DB=todo_database -e POSTGRES_USER=puser -e POSTGRES_PASSWORD=ppassword -d postgres && \
	docker ps | grep 'todo-postgres'
psql-stop:
	docker rm -f todo-postgres
migration-up:
	POSTGRESQL_URL='postgres://puser:ppassword@localhost:5432/todo_database?sslmode=disable' && \
    cd database && chmod +x ./migrate && ./migrate -database $$POSTGRESQL_URL -path ./migrations up
migration-down:
	POSTGRESQL_URL='postgres://puser:ppassword@localhost:5432/todo_database?sslmode=disable' && \
    cd database && chmod +x ./migrate && ./migrate  -database $$POSTGRESQL_URL -path ./migrations down -all
