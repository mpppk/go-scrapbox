REPO_OWNER = mpppk
REPO_NAME = go-scrapbox
ifdef update
  u=-u
endif

lint:
	gometalinter

test:
	go test ./...

circleci:
	circleci build -e GITHUB_TOKEN=$GITHUB_TOKEN

.PHONY: lint test circleci