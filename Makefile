BINARY=engine
test: 
	go test -v -cover -covermode=atomic ./...

engine:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BINARY} ./main.go

unittest:
	go test -short  ./...

test-coverage:
	go test -short -coverprofile cover.out -covermode=atomic ./...
	cat cover.out >> coverage.txt

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

docker:
	docker build -t balance-service .

run:
	docker-compose up -d

stop:
	docker-compose down

lint:
	golangci-lint run 

mock:
	moq -out ./user/repository/mock.go ./user/repository Querier

sqlc:
	sqlc generate

.PHONY: test engine unittest test-coverage clean docker run stop lint mock sqlc