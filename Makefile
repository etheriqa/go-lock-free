PACKAGES = github.com/etheriqa/go-lock-free/queue

.PHONY: all test bench

all: bench

test:
	go test -v $(PACKAGES)

bench:
	go test -cpu 1,2,4,8,16,32 -bench . -benchmem $(PACKAGES)
