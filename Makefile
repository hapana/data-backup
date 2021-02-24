build:
	go build -o bin/data-backup

smoke: build
	./bin/data-backup -p ./test/compose/
