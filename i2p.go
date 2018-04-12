package apti2p

import (
	"github.com/cryptix/goSam"
	"log"
	"net/http"
)

var (
	samClient *goSam.Client

	aptTransport *http.Transport
	aptClient    *http.Client
	err          error
)

func Init() {
	samClient, err = goSam.NewClient("127.0.0.1:7656")
	if Fatal(err, "SAM client created:", "127.0.0.1:7656") {
		aptTransport = &http.Transport{
			Dial: samClient.Dial,
		}
		aptClient = &http.Client{
			Transport: aptTransport,
		}
	}
}

func Fatal(err error, str ...string) bool {
	log.Println(str)
	if err != nil {
		log.Fatal(err, str)
		return false
	}
	return true
}
