COMMIT_HASH ?= $(shell git rev-parse HEAD 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X devteam/cmd.commitHash=$(COMMIT_HASH)"

.PHONY: build clean run

build:
	go build $(LDFLAGS) -o devteam .

clean:
	rm -f devteam

run:
	go run . $(ARGS)
