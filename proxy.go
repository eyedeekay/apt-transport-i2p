package apti2p

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type HttpConnection struct {
	Request  *http.Request
	Response *http.Response
}

type HttpConnectionChannel chan *HttpConnection

var connChannel = make(HttpConnectionChannel)

var addr = "127.0.0.1:7844"

func ProxyAddr() string {
	return addr
}

func ProxyHost() string {
	return strings.SplitN(addr, ":", 2)[0]
}

func ProxyPort() int {
	s := strings.Split(addr, ":")
	l := len(s) - 1
	i, err := strconv.Atoi(s[l])
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func SetProxyAddr(a, p string) string {
	addr = a + ":" + p
	return addr
}

type Proxy struct {
}

func NewProxy() *Proxy { return &Proxy{} }

func TranslateAddr(a string) string {
	t := strings.SplitN(a, ".i2p/", 2)
	x := strings.TrimPrefix(t[0], "i2p://")
	x = strings.TrimPrefix(x, "http://")
	x = strings.TrimPrefix(x, "https://")
	x = strings.TrimPrefix(x, "ftp://")
	p, err := Find(x)
	if err != nil {
		log.Fatal(err)
	}
	raddr := "http://" + ProxyHost() + ":" + p.TargetPort + "/" + t[1]
	log.Println(raddr)
	return raddr
}

func (p *Proxy) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	var req *http.Request

	client := &http.Client{}

	nr := TranslateAddr(r.RequestURI)
	req, err = http.NewRequest(r.Method, nr, r.Body)
	resp, err = client.Do(req)
	r.Body.Close()

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	for k, v := range resp.Header {
		wr.Header().Set(k, v[0])
	}
	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
	resp.Body.Close()

}

func InitClient() (*http.Client, error) {
	ret := http.Client{}
	return &ret, nil

}
