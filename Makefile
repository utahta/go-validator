
test:
	go test -race ./...

bench:
	go test -bench . -benchmem

pprof:
	go test -bench . -benchmem -cpuprofile cpu.out -memprofile mem.out

