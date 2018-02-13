run : postgres user-service
.PHONY: run
build:
	protoc -I. --go_out=plugins=micro:. \
		proto/user/auth.proto
image:
	docker build -t shippy-user-service .
postgres:
	docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres
	timeout 3 
user-service: 
	docker run -d --net="host" \
		-p 50051 \
		-e DB_HOST=192.168.99.100 \
		-e DB_PASSWORD=password \
		-e DB_NAME=postgres \
		-e DB_USER=postgres \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns \
		shippy-user-service