.PHONY: all fetchdeps build clean run start-wiremock

all: build
build: build/certgen
clean:
	rm -rf build

fetchdeps: go.mod Makefile
	go mod download

build/certgen: main/*.go Makefile
	go build -o build/certgen ./main
	CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w' -o build/certgen-static ./main

# `make run` will run the built app with a configuration suitable for the wiremock server below.
run: build/certgen Makefile
	REQUEST_SPEC_PATH=./test/resources/config.json \
	SERVICE_ACCOUNT_TOKEN_PATH=./test/resources/serviceAccountToken.jwt \
	K8S_AUTH_ROLE_NAME=mocked-cert-generator-role \
	VAULT_ADDR=http://localhost:80/ \
	KUBERNETES_CLUSTER_NAME=example-k8s-cluster \
	ENABLE_DEBUG=false \
	CONSOLE_FORMATTER=text \
		./build/certgen

# `make start-wiremock` will boot a wiremock server on port 80, listening for a login and a pki request.
# Pre-requisite of `make run`, which is configured to use this mock server.
start-wiremock:
	mkdir -p build
	cd build && wget --timestamping https://repo1.maven.org/maven2/com/github/tomakehurst/wiremock-standalone/2.24.0/wiremock-standalone-2.24.0.jar
	rm -rf build/mappings
	mkdir -p build/mappings
	cp ./test/resources/stubCertIssue.json build/mappings/
	cp ./test/resources/stubVaultLogin.json build/mappings/
	cd build && java -jar ./wiremock-standalone-2.24.0.jar --port 80 --disable-banner true --verbose true
