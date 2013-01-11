package main

import (
	"./magicport"
	"log"
)

func main() {
	cfg := new(magicport.Config)
	cfg.BindAddr = "127.0.0.1:4000"
	cfg.CertFile = "webcert.pem"
	cfg.KeyFile = "webkey.pem"
	cfg.Token = "#d423SDF!!@"

	log.Fatal(magicport.Run(cfg))
}
