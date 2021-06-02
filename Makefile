

EXEC_GO_CMD := docker-compose run --rm  go 
EXEC_SSH_CMD := docker-compose exec sshTestServer

build:
	rm -rf build && mkdir -p build && go build -o build/groac	
.PHONY: build

_test:
	go clean -testcache && go test -v ./...

test:
	$(EXEC_GO_CMD) make _test

prepare: # create required assets for testing
	docker-compose up -d
	$(EXEC_GO_CMD) rm -f /root/.groac/vms_key.pem /root/.groac/vms_key.pem.pub
	$(EXEC_GO_CMD) ssh-keygen -f /root/.groac/vms_key.pem -b 2048 -t rsa -q -N ""
	$(EXEC_SSH_CMD) cp /data/vms_key.pem.pub /config/.ssh/authorized_keys
clean:
	docker-compose down --volumes

.EXPORT_ALL_VARIABLES: