build:
	go build -o calc

run: build
	# time ./calc
	hyperfine --warmup 3 './calc'

profile:
	go tool pprof -http=localhost:8080 calc1.prof
