build:
	mkdir -p functions
	go get ./...
	go build -o functions/hello-lambda *.go

dev:
	go get ./...
	go build -o functions/hello-lambda *.go
	npm run start