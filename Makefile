.PHONY: build clean run

build:
	go build -o devteam .

clean:
	rm -f devteam

run:
	go run . $(ARGS)
