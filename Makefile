.PHONY: test build

NAME=gloc
BIN_PATH=./bin

build:
	mkdir -p $(BIN_PATH)
	go build -o $(BIN_PATH)/$(NAME) main.go result.go stack.go

clean:
	rm $(BIN_PATH)/*

# Run it like: bin/gloc --root=<dir where the code is> --ignore-test-files=<true/false> --exclude-dirs=<a,b> --exclude-files=<a.go,b.go>