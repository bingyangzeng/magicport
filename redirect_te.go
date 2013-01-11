package main

import (
	"./magicport"
	"log"
)

func main() {
	log.Printf("start")
	port := magicport.NewPort("tcp", "127.0.0.1:8080")
	port.AddInterface(magicport.NewRedirectInterface("127.0.0.1:8000"))
	port.ListenAndServe()
}
