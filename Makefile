build:
	go build -o calc

run: build
	./calc 2>&1 | less

time: build
	# time ./calc
	hyperfine --warmup 3 './calc'

test: build
	./calc | delta --max-line-length 0 --diff-so-fancy averages.txt -

profile: time
	go tool pprof -http=localhost:8080 calc1.prof
