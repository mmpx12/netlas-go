build:
	go build -ldflags="-w -s"

install:
	cp netlas /usr/bin/netlas

completion:
	go run main.go completion bash > /etc/bash_completion.d/netlas

all: build install

clean:
	rm -f netlas /usr/bin/netlas /etc/bash_completion.d/netlas
