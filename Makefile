install:
	go get github.com/tools/godep
	go get github.com/alecthomas/gometalinter
	gometalinter --install --update --force
	go install

lint:
	gometalinter --deadline=60s --vendor --cyclo-over=20 --disable=vetshadow,dupl ./...

test:
	go test $(shell go list ./... | grep -v '/vendor/') --cover

package:
	go clean
	OS="darwin"
	CGO_ENABLED=0 GOOS=$$OS go build -a -installsuffix cgo -o oauth2_central .
	./oauth2_central -version

all: install lint test package
