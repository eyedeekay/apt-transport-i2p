VERSION = 0.1

lib:
	go build .

build:
	go build -o ./bin/apt-transport-i2p ./apt-transport-i2p

release:
	GOOS=linux GOARCH=amd64 go build \
		-a \
		-tags netgo \
		-ldflags '-w -extldflags "-static"' \
		-o ./bin/apt-transport-i2p \
		./apt-transport-i2p

orig:
	tar --exclude=.git --exclude=bin --exclude=debian -czvf ../apt-transport-i2p_0.1.orig.tar.gz .

install:
	mkdir -p /etc/apt-transport-i2p/
	install -m755 bin/apt-transport-i2p /usr/lib/apt/methods/i2psam
	install etc/apt-transport-i2p/apt-transport-i2p.conf /etc/apt-transport-i2p/apt-transport-i2p.conf

