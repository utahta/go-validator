
test:
	go test -race ./...

bench:
	go test -bench . -benchmem

