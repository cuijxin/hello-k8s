all: gotool
	@go build -v -o atom-server .
clean:
	rm -f atom-server
	rm -rf docs
gotool:
	gofmt -w .
help:
	@echo "make - compile the source code"
	@echo "make clean - remove binary file and vim swp files"
	@echo "make gotool - run go tool 'fmt' and 'vet'"

.PHONY: clean gotool help
