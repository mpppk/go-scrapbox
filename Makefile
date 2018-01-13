REPO_OWNER = mpppk
REPO_NAME = go-scrapbox
ifdef update
  u=-u
endif

deps:
	dep ensure

setup:
	go get ${u} github.com/golang/dep/cmd/dep
	go get ${u} gopkg.in/alecthomas/gometalinter.v2

lint: deps
	gometalinter

test: deps
	go test ./...

circleci:
	circleci build

.PHONY: deps setup lint test circleci