ME := go-shortener
SRC := $(shell find . -name "*.go" -type f)
REGISTRY := registry.camderwin.us
REPO := $(REGISTRY)/$(ME)
SERVICE_NAME := app

build:
	docker-compose build

run:
	docker-compose run --rm --service-ports $(SERVICE_NAME)

kill:
	CONTAINERS=$$(docker ps --format "{{.Names}}" | grep $(ME)_$(SERVICE_NAME)) && \
	for CONTAINER in $$CONTAINERS ; do docker rm -f $$CONTAINER ; done

bash:
	docker-compose run --rm $(SERVICE_NAME) bash

savedeps:
	docker-compose run --rm $(SERVICE_NAME) godep save ./...

fmt: .fmt.ts

.fmt.ts:
	docker-compose run --rm $(SERVICE_NAME) gofmt -l -s -w .

lint:
	@docker-compose run --rm $(SERVICE_NAME) golint

test: .tests.ts

.tests.ts: $(SRC)
	docker-compose run --rm $(SERVICE_NAME) go test -cover -v && \
	touch $@
