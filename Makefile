all: test

test:
	luac test.lua
	go test -v ./


