build:
	go build -o calc

run: build
	time ./calc
