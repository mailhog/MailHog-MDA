DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: deps fmt combined

combined:
	go install .

release: release-deps
	gox -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" .

fmt:
	go fmt ./...

deps:

test-deps:
	go get github.com/smartystreets/goconvey

release-deps:
	go get github.com/mitchellh/gox

.PNONY: all combined release fmt deps test-deps release-deps
