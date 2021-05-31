

build:
	rm -rf build && mkdir -p build && go build -o build/groac	
.PHONY: build

test:
	go clean -testcache && go test -v ./...


.EXPORT_ALL_VARIABLES: