build:
	go build -o calc

run: build
	# time ./calc
	hyperfine --warmup 3 './calc'

test: build
	./calc | delta --max-line-length 0 --diff-so-fancy averages.txt -

profile:
	go tool pprof -http=localhost:8080 calc1.prof
