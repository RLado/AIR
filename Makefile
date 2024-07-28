BINARY_NAME:=oink
DESTDIR:= 
 
build:
	@go build -ldflags "-s -w" -o ${BINARY_NAME} src/main.go

clean:
	@go clean
	@rm ${BINARY_NAME}