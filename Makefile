VERSION=${shell cat VERSION}
CMD=cmd/mdqi/mdqi
PROJECT_USERNAME=handlename
PROJECT_REPONAME=mdqi

$(CMD): *.go */**/*.go
	go build -v -ldflags '-X main.version=$(VERSION)' -o $@ cmd/mdqi/main.go

.PHONY: dist
dist:
	cd cmd/mdqi && gox -os="linux darwin" -arch="386 amd64" -output "../../dist/mdqi-$(VERSION)-{{.OS}}-{{.Arch}}" -ldflags '-X main.version=$(VERSION)'

.PHONY: publish
publish:
	ghr -u '$(PROJECT_USERNAME)' -r '$(PROJECT_REPONAME)' --replace '$(VERSION)' dist/

.PHONY: clean
clean:
	-rm -rf $(CMD) dist/*
