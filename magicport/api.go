package magicport

import (
	"log"
)

func APIInit(app *App) {
	app.AddHandle("/port/new/redirect/", handleNewRedirect)
}

func handleNewRedirect(c *Context) {
	/*w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))*/
	c.WriteString("hello")

	log.Printf("url: %s", c.Request.URL)
}
