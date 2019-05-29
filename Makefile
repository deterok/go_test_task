DOCKER_DIR = docker

.PHONY: dc-%
dc-%:
	$(DOCKER_DIR)/$*.sh

.PHONY: up-build
up-build:
	$(DOCKER_DIR)/up.sh --build

.PHONY: up
up: dc-up

.PHONY: down
down:
	$(DOCKER_DIR)/down.sh -v

.PHONY: start
start: dc-start

.PHONY: stop
stop: dc-stop

test:
	$(DOCKER_DIR)/run.sh payments go test  -timeout 60s -v  ./...
