package meilo

import (
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-smtp"
)

type server struct {
	Port     string
	Password string
	User     string
	Host     string

	senderOpts []senderOption
}

func (bkd *server) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &session{
		username: bkd.User,
		password: bkd.Password,
		sender:   newSender(bkd.senderOpts...),
	}, nil
}

func (bkd *server) Addr() string {
	return bkd.Host + ":" + bkd.Port
}

// Start starts the SMTP server with the given options.
func (bkd *server) run() error {
	stp := smtp.NewServer(bkd)
	stp.Addr = bkd.Host + ":" + bkd.Port
	stp.Domain = bkd.Host
	stp.WriteTimeout = 10 * time.Second
	stp.ReadTimeout = 10 * time.Second
	stp.MaxMessageBytes = 1024 * 1024
	stp.MaxRecipients = 50
	stp.AllowInsecureAuth = true

	log.Println("Starting SMTP server at", stp.Addr)
	if err := stp.ListenAndServe(); err != nil {
		return fmt.Errorf("meilo: failed to start server: %v", err)
	}

	return nil
}
