package main

import (
	. ".."
)

func main() {
	//Init("/etc/apt-transport-i2p/apt-transport-i2p.conf")
    Init()
	a := AptMethod{}
	a.Run()
}
