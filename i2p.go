package apti2p

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/eyedeekay/portcheck"
	"github.com/eyedeekay/sam-forwarder"
	"github.com/eyedeekay/sam-forwarder/config"
)

var (
	aptTunnel []*samforwarder.SAMClientForwarder
	err       error
)

func ReadConfigs() []string {
	var s []string
	files, err := ioutil.ReadDir("/etc/apt/sources.list.d/")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		s = append(s, readConfigFile(file.Name())...)
	}
	log.Println(s)
	return s
}

func readConfigFile(filename string) []string {
	file, err := ioutil.ReadFile("/etc/apt/sources.list.d/"+filename)
	var s []string
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(file), "\n")
	for _, l := range lines {
		if !strings.HasPrefix(l, "#") {
			if strings.Contains(l, "i2p://") {
				t := strings.TrimPrefix(l, "i2p://")
				t = strings.TrimPrefix(t, "http://")
				t = strings.TrimPrefix(t, "https://")
				r := strings.SplitN(t, ".i2p", 2)
				s = append(s, r[0])
			}
		}
	}
	return s
}

func Find(addr string) (*samforwarder.SAMClientForwarder, error) {
	for _, a := range aptTunnel {
		if addr == a.Destination() {
			return a, nil
		}
	}
	return nil, nil
}

func Wait(tmp *samforwarder.SAMClientForwarder) {
	for {
		if len(tmp.Base32()) > 51 {
			log.Println("base32: ", tmp.Base32())
			break
		} else {
			log.Println("waiting for address")
		}
	}
}

func Init(addr string) (*samforwarder.SAMClientForwarder, error) {
	globalConf, err := i2ptunconf.NewI2PTunConf("/usr/share/apt-transport-i2p/apt.ini")
	if err != nil {
		return nil, err
	}
    log.Println("Initializing tunnel")
	SetProxyAddr(globalConf.TargetHost, globalConf.TargetPort)
	if client, err := Find(addr); client != nil && err == nil {
		return client, nil
	} else if err != nil {
		return nil, err
	}

	globalConf.ClientDest = addr
	i, err := strconv.Atoi(globalConf.TargetPort)
	if err != nil {
		return nil, err
	}
	globalConf.TargetPort = pc.SFL(i)
	tmp, err := i2ptunconf.NewSAMClientForwarderFromConf(globalConf)
	if err != nil {
		return nil, err
	}
	aptTunnel := append(aptTunnel, tmp)
	end := len(aptTunnel) - 1
	go aptTunnel[end].Serve()

	Wait(aptTunnel[end])

	return aptTunnel[end], nil
}
