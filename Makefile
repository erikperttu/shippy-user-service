run : postgres user-service
stop: stop-postgres stop-user-service
.PHONY: run
.PHONY: stop
build:
	protoc -I. --go_out=plugins=micro:. \
		proto/user/auth.proto
image:
	docker build -t shippy-user-service .
postgres:
	docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres
	timeout 3 
user-service: 
	docker run --net="host" \
		-p 50051 \
		-e DB_HOST=192.168.99.100 \
		-e DB_PASSWORD=password \
		-e DB_NAME=postgres \
		-e DB_USER=postgres \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns \
		shippy-user-service
stop-postgres:
	docker ps -q --filter ancestor=postgres | xargs -r docker stop
stop-user-service:
	docker ps -q --filter ancestor=shippy-user-service | xargs -r docker stop