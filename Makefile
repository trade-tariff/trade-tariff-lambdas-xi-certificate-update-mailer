.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/handler xi-certificate-update-mailer/main.go

clean:
	rm -rf ./bin xi-certificate-update-mailer/bin

deploy: clean build
	cd xi-certificate-update-mailer && sls deploy --verbose

test:
	cd xi-certificate-update-mailer && go test -v ./...
