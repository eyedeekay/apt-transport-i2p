package apti2p

import (
	"crypto/tls"
	"github.com/eyedeekay/gosam"
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

func Init(conf ...string) {
	addr, port, il, ol, iq, oq, biq, boq, debug := ParseConfig(conf)
	samClient, err = goSam.NewClientFromOptions(goSam.SetHost(addr), goSam.SetPort(port), goSam.SetInLength(il), goSam.SetOutLength(ol), goSam.SetInQuantity(iq), goSam.SetOutQuantity(oq), goSam.SetInBackups(biq), goSam.SetOutBackups(boq), goSam.SetDebug(debug))
	if Fatal(err, "SAM client created:", "127.0.0.1:7656") {
		aptTransport = &http.Transport{
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   10,
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: time.Duration(2) * time.Minute,
			ExpectContinueTimeout: time.Duration(2) * time.Minute,
			IdleConnTimeout:       time.Duration(6) * time.Minute,
			TLSNextProto:          make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			Dial:                  samClient.Dial,
		}
		aptClient = &http.Client{
			Timeout:       time.Duration(6) * time.Minute,
			Transport:     aptTransport,
			Jar:           nil,
			CheckRedirect: nil,
		}
	}
}

func ParseConfig(conf []string) (string, string, uint, uint, uint, uint, uint, uint, bool) {
	samhost, samport := "127.0.0.1", "7656"
	inlen, outlen, inquantity, outquantity, backupin, backupout := uint(2), uint(2), uint(15), uint(2), uint(5), uint(2)
    debug := false
    if len(conf) == 1 {
        //TODO: Read in config file and put the variables into the return values. for/range+switch
    }
	return samhost, samport, inlen, outlen, inquantity, outquantity, backupin, backupout, debug
}

func Fatal(err error, str ...string) bool {
	log.Println(str)
	if err != nil {
		log.Fatal(err, str)
		return false
	}
	return true
}
