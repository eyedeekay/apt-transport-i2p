package apti2p

import (
	"github.com/cryptix/goSam"
	"log"
	"net/http"
	"time"
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
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   4,
			DisableKeepAlives:     false,
			ResponseHeaderTimeout: time.Duration(2) * time.Minute,
			ExpectContinueTimeout: time.Duration(2) * time.Minute,
			IdleConnTimeout:       time.Duration(6) * time.Minute,
            //TLSNextProto:          make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			Dial:                  samClient.Dial,
		}
		aptClient = &http.Client{
			Timeout:   time.Duration(6) * time.Minute,
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
