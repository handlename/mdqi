VERSION=${shell cat VERSION}
CMD=dist/mdqi-$(VERSION)

$(CMD): */**/*.go
	go build -v -ldflags '-X main.version=$(VERSION)' -o $@ cmd/mdqi/main.go

.PHONY: clean
clean:
	-rm $(CMD)
