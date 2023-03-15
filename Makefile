build:
	@go build ./cmd/farmer/main.go

run-w:
	@./main.exe

run-l:
	@./main

all-w: build run-w
all-l: build run-l