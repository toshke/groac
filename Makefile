

RUN_GO_CMD := docker-compose run --rm  go 

build:
	rm -rf build && mkdir -p build && go build -o build/groac	
.PHONY: build

_test:
	go clean -testcache && go test -v ./...

test:
	$(RUN_GO_CMD) make _test


clean:
	docker-compose down --volumes

.EXPORT_ALL_VARIABLES: