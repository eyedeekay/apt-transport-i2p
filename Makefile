
lib:
	go build .


build:
	cd main && \
		go build -o ../bin/apt-transport-i2p

install:
	install -m755 bin/apt-transport-i2p /usr/lib/apt/methods/i2p
