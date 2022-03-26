build:
	go build -ldflags="-w -s"

install:
	cp netlas /usr/bin/netlas

all: build install

clean:
	rm -f netlas /usr/bin/netlas
