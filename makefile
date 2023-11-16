BINARY_NAME=gophermart

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin ./cmd/gophermart/
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux ./cmd/gophermart/
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows.exe ./cmd/gophermart/

run: build
	./${BINARY_NAME}


clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows.exe

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet ./...

lint:
	go fmt ./...