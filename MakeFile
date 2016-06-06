install:
	go get github.com/tools/godep
	go get github.com/alecthomas/gometalinter
	gometalinter --install --update --force
	go install

lint:
	gometalinter --deadline=40s --vendor --cyclo-over=20 --disable=vetshadow ./...

test:
	go test github.com/vedhavyas/oauth2_central/... --cover

package:
	go clean
	CGO_ENABLED=0 GOOS=$(shell echo `uname` | awk '{print tolower($0)}') go build -a -installsuffix cgo -o oauth2_central .

all: install lint test package
