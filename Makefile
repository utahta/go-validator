
test:
	go test -race ./...

coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic

bench:
	go test -bench . -benchmem

pprof:
	go test -bench . -benchmem -cpuprofile cpu.out -memprofile mem.out

changelog:
	git-chglog -o CHANGELOG.md

godoc:
	godoc -http=:6060

