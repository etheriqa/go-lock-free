PACKAGES = github.com/etheriqa/go-lock-free/queue

.PHONY: bench

bench:
	go test -cpu 1,2,4,8,16,32 -bench . -benchmem $(PACKAGES) 
