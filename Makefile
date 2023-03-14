build:
	@go build main.go

run:
	@./main

clean:
	@go clean .

all: build run clean