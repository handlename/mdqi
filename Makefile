VERSION=${shell cat VERSION}
CMD=dist/mdqi-$(VERSION)
PROJECT_USERNAME=
PROJECT_REPONAME=

$(CMD): */**/*.go
	go build -v -ldflags '-X main.version=$(VERSION)' -o $@ cmd/mdqi/main.go

.PHONY: dist
dist:
	cd cmd/mdqi && gox -os="linux darwin" -arch="386 amd64" -output "../../dist/mdqi-$(VERSION)-{{.OS}}-{{.Arch}}/{{.Dir}}" -ldflags '-X main.version=$(VERSION)'

.PHONY: publish
publish:
	ghr -u $PROJECT_USERNAME -r $ROJECT_REPONAME --replace $(cat VERSION) dist/

.PHONY: clean
clean:
	-rm -rf dist/mdqi-*
