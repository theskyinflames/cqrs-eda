default: test

test:
	go test -v -race ./...

lint:
	revive -config ./revive.toml
	go mod tidy -v && git --no-pager diff --quiet go.mod go.sum

tools:  tool-moq tool-revive

tool-revive:
	go install github.com/mgechev/revive@main

tool-moq:
	go install github.com/matryer/moq@main

todo:
	find . -name '*.go' \! -name '*_generated.go' -prune | xargs grep -n TODO


