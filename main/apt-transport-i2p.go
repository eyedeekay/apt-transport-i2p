package main

import (
    "os"
	//. "github.com/eyedeekay/apt-transport-i2p"
    . ".."
)

func main() {
    configfile := os.Getenv("APT_TRANSPORT_SAM_CONFIG")
	Init(configfile)
	a := AptMethod{}
	a.Run()
}
