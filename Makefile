
RUN_GOLANG := docker-compose run -w /src --rm golang
PROJECT_ROOT := $(shell pwd)

_build:
	bash tooling/scripts/build.sh
	
build:
	$(RUN_GOLANG) make _build
.PHONY: build

shell:
	$(RUN_GOLANG) bash

.EXPORT_ALL_VARIABLES: