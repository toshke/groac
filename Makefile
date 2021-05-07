

build:
	rm -rf build && mkdir -p build && go build -o build/groac	
.PHONY: build

test:
	go test ./...


.EXPORT_ALL_VARIABLES: