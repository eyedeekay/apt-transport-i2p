package apti2p

import (
	"crypto/tls"
	"github.com/eyedeekay/gosam"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	samClient *goSam.Client

	aptTransport *http.Transport
	aptClient    *http.Client
	err          error
)

func Init(conf ...string) {
	addr, port, il, ol, iq, oq, biq, boq, gzip, debug := ParseConfig(conf)
	samClient, err = goSam.NewClientFromOptions(
		goSam.SetHost(addr),
		goSam.SetPort(port),
		goSam.SetInLength(il),
		goSam.SetOutLength(ol),
		goSam.SetInQuantity(iq),
		goSam.SetOutQuantity(oq),
		goSam.SetInBackups(biq),
		goSam.SetOutBackups(boq),
		goSam.SetCompression(gzip),
		goSam.SetDebug(debug),
        goSam.SetUnpublished(true),
	)
	if Fatal(err, "SAM client created:", "127.0.0.1:7656") {
		aptTransport = &http.Transport{
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   2,
			DisableKeepAlives:     false,
			ResponseHeaderTimeout: time.Duration(10) * time.Minute,
			ExpectContinueTimeout: time.Duration(10) * time.Minute,
			IdleConnTimeout:       time.Duration(10) * time.Minute,
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

func ParseConfig(conf []string) (string, string, uint, uint, uint, uint, uint, uint, bool, bool) {
	samhost, samport := "127.0.0.1", "7656"
	inlen, outlen, inquantity, outquantity, backupin, backupout := uint(2), uint(2), uint(2), uint(2), uint(1), uint(1)
	debug, gzip := false, true
	if len(conf) == 1 {
		if conf[0] != "" {
			confstring, err := ioutil.ReadFile(conf[0])
			if err != nil {
				log.Fatal(err)
			}
			items := strings.Split(string(confstring), "\n")
			for _, item := range items {
				kv := strings.Split(item, "=")
				if len(kv) == 2 {
					if strings.HasPrefix(kv[0], "samhost") {
						samhost = kv[1]
					}
					if strings.HasPrefix(kv[0], "samport") {
						samport = kv[1]
					}
					if strings.HasPrefix(kv[0], "inlen") {
						val, err := strconv.Atoi(kv[1])
						if err == nil {
							inlen = uint(val)
						}
					}
					if strings.HasPrefix(kv[0], "outlen") {
						val, err := strconv.Atoi(kv[1])
						if err == nil {
							outlen = uint(val)
						}
					}
					if strings.HasPrefix(kv[0], "inquantity") {
						val, err := strconv.Atoi(kv[1])
						if err == nil {
							inquantity = uint(val)
						}
					}
					if strings.HasPrefix(kv[0], "outquantity") {
						val, err := strconv.Atoi(kv[1])
						if err == nil {
							outquantity = uint(val)
						}
					}
					if strings.HasPrefix(kv[0], "backupin") {
						val, err := strconv.Atoi(kv[1])
						if err == nil {
							backupin = uint(val)
						}
					}
					if strings.HasPrefix(kv[0], "backupout") {
						val, err := strconv.Atoi(kv[1])
						if err == nil {
							backupout = uint(val)
						}
					}
					if strings.HasPrefix(kv[0], "gzip") {
						val, err := strconv.ParseBool(kv[1])
						if err == nil {
							gzip = val
						}
					}
					if strings.HasPrefix(kv[0], "debug") {
						val, err := strconv.ParseBool(kv[1])
						if err == nil {
							debug = val
						}
					}
				}
			}
		}
	}
	return samhost, samport, inlen, outlen, inquantity, outquantity, backupin, backupout, gzip, debug
}

func Fatal(err error, str ...string) bool {
	log.Println(str)
	if err != nil {
		log.Fatal(err, str)
		return false
	}
	return true
}
