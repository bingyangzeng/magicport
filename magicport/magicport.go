package magicport

import (
	"log"
)

type Config struct {
	BindAddr string
	KeyFile  string
	CertFile string
	Token    string
}

func Run(config *Config) error {
	log.Printf("Start at https://%s", config.BindAddr)

	app := NewApp("127.0.0.1")
	APIInit(app)

	err := app.ListenAndServeTLS(config.BindAddr, config.CertFile, config.KeyFile)
	return err
}
