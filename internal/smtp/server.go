package smtp

import (
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/wawandco/meilo/internal/models"
)

type Server struct {
	Port, Password, User, Host string
	WebPort                    string
	SaveFn                     func(models.Email)
}

func (bkd *Server) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &session{
		username: bkd.User,
		password: bkd.Password,
		saveFn:   bkd.SaveFn,
	}, nil
}

func (bkd *Server) Addr() string {
	return bkd.Host + ":" + bkd.Port
}

// Run starts the SMTP server with the given options.
func (bkd *Server) Run() error {
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
