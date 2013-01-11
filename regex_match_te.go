package main

import "./magicport"

func main() {
	port := magicport.NewPort("tcp", "127.0.0.1:8080")
	inter := magicport.NewRegexMatchInterface(
		"220.166.52.189:80", "^(GET|POST) /oj/", []byte("\r\n"), 0)
	port.AddInterface(inter)
	port.AddInterface(magicport.NewRedirectInterface("127.0.0.1:8000"))
	port.ListenAndServe()
}
