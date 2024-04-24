.PHONY: build clean deploy-development deploy-staging deploy-production test lint install-linter

build:
	cd xi-certificate-update-mailer && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap

clean:
	rm -rf boostrap

deploy-development: clean build
	STAGE=development serverless deploy --verbose

deploy-staging: clean build
	STAGE=staging serverless deploy --verbose

deploy-production: clean build
	STAGE=production serverless deploy --verbose

test: clean build
	cd xi-certificate-update-mailer && go test ./...

lint:
	cd xi-certificate-update-mailer && golangci-lint run
