build:
	go build -o bin/data-backup

smoke: build
	cd test && ../bin/data-backup
