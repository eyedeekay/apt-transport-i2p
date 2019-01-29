#!/usr/bin/make -f

## Copyright (C) 2012 - 2018 ENCRYPTED SUPPORT LP <adrelanos@riseup.net>
## See the file COPYING for copying conditions.

## genmkfile - Makefile - version 1.5

## This is a copy.
## master location:
## https://github.com/Whonix/genmkfile/blob/master/usr/share/genmkfile/Makefile

fmt:
	gofmt -w *.go */*.go

GENMKFILE_PATH ?= /usr/share/genmkfile
GENMKFILE_ROOT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

export GENMKFILE_PATH
export GENMKFILE_ROOT_DIR

include $(GENMKFILE_PATH)/makefile-full


VERSION = 0.1

lib:
	go build .

deps:
	go get -u github.com/eyedeekay/littleboss
	go get -u github.com/eyedeekay/portcheck
	go get -u github.com/eyedeekay/sam-forwarder

build:
	go build -a \
		-tags netgo \
		-o usr/bin/apt-transport-i2p \
		./main

release:
	go build -a -tags netgo \
		-ldflags '-w -extldflags "-static"' \
		-o usr/bin/apt-transport-i2p \
		./main

config:
	rm -rfv usr/share/apt-transport-i2p usr/lib/apt/methods
	mkdir -pv usr/share/apt-transport-i2p usr/lib/apt/methods
	@echo '#! /bin/sh' | tee usr/lib/apt/methods/i2p
	@echo '/usr/bin/apt-transport-i2p --littleboss=status 2>/dev/null || \' | tee -a usr/lib/apt/methods/i2p
	@echo '    nohup /usr/bin/apt-transport-i2p --littleboss=start & exit' | tee -a usr/lib/apt/methods/i2p
	@echo '/usr/bin/apt-transport-i2p --littleboss=status 2>/dev/null && \' | tee -a usr/lib/apt/methods/i2p
	@echo '    /usr/bin/apt-transport-i2p $$@ &' | tee -a usr/lib/apt/methods/i2p
	chmod +x usr/lib/apt/methods/i2p
	ln -sf usr/lib/apt/methods/i2p usr/lib/apt/methods/i2p+http
	ln -sf usr/lib/apt/methods/i2p usr/lib/apt/methods/i2p+https
	ln -sf usr/lib/apt/methods/i2p usr/lib/apt/methods/i2p+ftp
	@echo 'type = client' | tee usr/share/apt-transport-i2p/apt.ini
	@echo 'host = 127.0.0.1' | tee -a usr/share/apt-transport-i2p/apt.ini
	@echo 'port = 7844' | tee -a usr/share/apt-transport-i2p/apt.ini
	@echo 'inbound.length = 2' | tee -a usr/share/apt-transport-i2p/apt.ini
	@echo 'outbound.length = 2' | tee -a usr/share/apt-transport-i2p/apt.ini
	@echo 'inbound.quantity = 1' | tee -a usr/share/apt-transport-i2p/apt.ini
	@echo 'outbound.quantity = 1' | tee -a usr/share/apt-transport-i2p/apt.ini
	@echo 'inbound.backupQuantity = 1' | tee -a usr/share/apt-transport-i2p/apt.ini
	@echo 'outbound.backupQuantity = 1' | tee -a usr/share/apt-transport-i2p/apt.ini

description-pak:

checkinstall: release description-pak
	checkinstall --default \
		--install=no \
		--fstrans=yes \
		--maintainer=hankhill19580@gmail.com \
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

test:
	./usr/lib/apt/methods/i2p -littleboss=start &
	@echo "started self-supervising server"
	sleep 5
	./usr/lib/apt/methods/i2p


kill:
	./usr/lib/apt/methods/i2p -littleboss=stop
