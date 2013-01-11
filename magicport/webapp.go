package magicport

import (
	// "log"
	"net/http"
)

type App struct {
	domain   string
	handlers map[string]HandlerFunc
	sessions map[string]*map[string]string
}

type Context struct {
	app      *App
	Request  *http.Request
	Response *http.ResponseWriter
	Sessions *map[string]string
}

type HandlerFunc func(*Context)

func (self *Context) WriteString(data string) (int, error) {
	return (*self.Response).Write([]byte(data))
}

func (self *Context) SetCookie(name, value string) {
	cookie := http.Cookie{Name: name, Value: value, Path: "/",
		MaxAge: 0, Secure: true, HttpOnly: false}
	http.SetCookie(*self.Response, &cookie)
}

func (self *Context) Cookie(name string) (string, error) {
	cookie, err := self.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func NewApp(domain string) *App {
	app := new(App)
	app.domain = domain
	app.handlers = make(map[string]HandlerFunc)
	app.sessions = make(map[string]*map[string]string)
	return app
}

func (self *App) makeContext(w *http.ResponseWriter, req *http.Request) *Context {
	ctx := new(Context)
	ctx.Request = req
	ctx.Response = w
	ctx.app = self

	sessionId, err := ctx.Cookie("session")
	if err != nil {
		sessionId = NewSessionId()
		session := make(map[string]string)
		self.sessions[sessionId] = &session
		ctx.Sessions = &session
		ctx.SetCookie("session", sessionId)
	} else {
		ctx.Sessions = self.sessions[sessionId]
	}

	return ctx
}

func (self *App) ListenAndServeTLS(addr, cert, key string) error {
	return http.ListenAndServeTLS(addr, cert, key, self)
}

func (self *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := self.makeContext(&w, req)
	handler, ok := self.handlers[req.URL.Path]

	if ok {
		handler(ctx)
	} else {
		http.NotFound(w, req)
	}
}

func (self *App) AddHandle(url string, handler HandlerFunc) {
	self.handlers[url] = handler
}
