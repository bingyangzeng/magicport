package main

import (
	"./magicport"
)

func main() {
	port := magicport.NewPort("tcp", "127.0.0.1:9000")
	port.AddInterface(magicport.NewAnyPortInterface([]byte{}, []byte("pwd")))
	port.ListenAndServe()
}
