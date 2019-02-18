package main

import (
	. ".."
	"context"
	"github.com/eyedeekay/littleboss"
	"github.com/eyedeekay/portcheck"
	"log"
	"net/http"
)

func main() {
	lb := littleboss.New("apt-transport-i2p")
	lb.Persist = true
	lb.Run(func(ctx context.Context) {
		if b, e := pc.CheckLocal(ProxyPort(), ProxyHost()); b {
			for _, i2paddr := range ReadConfigs() {
				_, err := Init(i2paddr)
				if err != nil {
					log.Fatal(err)
				}
			}
			handler := NewProxy()
			log.Println("Starting proxy server on", ProxyAddr())
			if err := http.ListenAndServe(ProxyAddr(), handler); err != nil {
				log.Fatal("ListenAndServe:", err)
			}
		} else if e != nil {
			log.Fatal(e)
		}
		log.Println("Tunnels are running.")
		a := AptMethod{}
		a.Run()
	})
}
