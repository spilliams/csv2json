.PHONY: run build install
run:
	go run src/cmd/csv2json/main.go

build:
	go build -o bin/csv2json src/cmd/csv2json/main.go

install:
	go build -o $$GOPATH/bin/csv2json src/cmd/csv2json/main.go
