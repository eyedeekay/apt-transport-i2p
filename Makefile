VERSION = 0.2

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
