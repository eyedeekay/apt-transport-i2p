VERSION = 0.1

lib:
	go build .

build:
	cd main && \
		go build -o ../bin/apt-transport-i2p

release:
	cd main && \
	GOOS=linux GOARCH=amd64 go build \
		-a \
		-tags netgo \
		-ldflags '-w -extldflags "-static"' \
		-o bin/apt-transport-i2p

install:
	mkdir -p /etc/apt-transport-i2p/
	install -m755 bin/apt-transport-i2p /usr/lib/apt/methods/i2p
	install etc/apt-transport-i2p/apt-transport-i2p.conf /etc/apt-transport-i2p/apt-transport-i2p.conf

description-pak:

checkinstall: release description-pak
	checkinstall --default \
		--install=no \
		--fstrans=yes \
		--maintainer=problemsolver@openmailbox.org \
		--pkgname="apt-transport-i2p" \
		--pkgversion="$(VERSION)" \
		--pkglicense=gpl \
		--pkggroup=net \
		--pkgsource=./ \
		--deldoc=yes \
		--deldesc=yes \
		--delspec=yes \
		--backup=no \
		--pakdir=../

checkinstall-arm: build-arm description-pak static-include static-exclude
	checkinstall --default \
		--install=no \
		--fstrans=yes \
		--maintainer=problemsolver@openmailbox.org \
		--pkgname="apt-transport-i2p" \
		--pkgversion="$(VERSION)-arm" \
		--pkglicense=gpl \
		--pkggroup=net \
		--pkgsource=./ \
		--deldoc=yes \
		--deldesc=yes \
		--delspec=yes \
		--backup=no \
		--exclude=arm-exclude \
		--include=arm-include \
		--pakdir=../
