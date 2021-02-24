build:
	go build -o bin/data-backup

smoke: build
	CURRENT_UID=$(shell id -u):$(shell id -g) \
		./bin/data-backup -p ./test/compose/
