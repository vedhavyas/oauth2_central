install:
	go get github.com/tools/godep
	go get github.com/alecthomas/gometalinter
	gometalinter --install --update --force
	go install

lint:
	gometalinter --deadline=60s --vendor --cyclo-over=20 --disable=vetshadow ./...

test:
	go test $(go list ./... | grep -v '/vendor/') --cover

package:
	go clean
	CGO_ENABLED=0 GOOS=$(shell echo `uname` | awk '{print tolower($0)}') go build -a -installsuffix cgo -o oauth2_central .
	./oauth2_central -version

all: install lint test package
