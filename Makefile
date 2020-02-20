VERSION = 0.2

GO111MODULE=off

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

link:
	rm -rf $(HOME)/go/src/github.com/eyedeekay/gosam
	ln -sf $(HOME)/go/src/github.com/eyedeekay/goSam $(HOME)/go/src/github.com/eyedeekay/gosam