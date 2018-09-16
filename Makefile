
test:
	go test -race ./...

bench:
	go test -bench . -benchmem

prof:
	go test -bench . -benchmem -cpuprofile cpu.out -memprofile mem.out

