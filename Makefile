.PHONY: build clean deploy

build:
	cd xi-certificate-update-mailer && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../bin/handler

clean:
	rm -rf ./bin xi-certificate-update-mailer/bin

deploy: clean build
	cd xi-certificate-update-mailer && sls deploy --verbose

test:
	cd xi-certificate-update-mailer && go test ./...
