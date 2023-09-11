.PHONY: build clean deploy-development deploy-staging deploy-production test

build:
	cd xi-certificate-update-mailer && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../bin/handler

clean:
	rm -rf ./bin

deploy-development: clean build
	STAGE=development serverless deploy --verbose

deploy-staging: clean build
	STAGE=staging serverless deploy --verbose

deploy-production: clean build
	STAGE=production serverless deploy --verbose

test: clean build
	cd xi-certificate-update-mailer && go test ./...
