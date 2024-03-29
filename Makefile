server:
	cd service && go get ./...
	cd service && go build -o server main.go

client:
	cd api && go get ./...
	cd api && go build -o client main.go

run: server client
	service/server &
	api/client &
	echo "Server and Client started successfully, Listening at [::1]:55055"
