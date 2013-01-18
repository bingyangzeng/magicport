package main

import (
	"./magicport"
	"flag"
	"fmt"
)

var bind = flag.String("bind", "", "address to bind")
var key = flag.String("key", "", "auth key")

func main() {
	flag.Parse()

	port := magicport.NewPort("tcp", *bind)
	fmt.Println(*key)
	port.AddInterface(magicport.NewAnyPortInterface([]byte{}, []byte(*key)))
	port.ListenAndServe()
}
