cli:
	go build -mod vendor -o bin/server cmd/server/main.go

lambda-handlers:
	@make lambda-server

lambda-server:	
	if test -f main; then rm -f main; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOOS=linux go build -mod vendor -o main cmd/server/main.go
	zip server.zip main
	rm -f main
