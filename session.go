package meilo

import (
	"errors"
	"io"
	"log"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

var e = email{}

// A Session is returned after successful login.
type session struct {
	username string
	password string
}

// AuthMechanisms returns a slice of available auth mechanisms; only PLAIN is supported.
func (s *session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

// Auth is the handler for supported authenticators.
func (s *session) Auth(mech string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(identity, username, password string) error {
		if username != s.username || password != s.password {
			return errors.New("invalid username or password")
		}
		return nil
	}), nil
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	e.From = from
	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	e.To = append(e.To, to)
	return nil
}

func (s *session) Data(r io.Reader) error {
	_, err := io.Copy(&e.Body, r)
	if err != nil {
		return err
	}

	return nil
}

func (s *session) Logout() error { return nil }

func (s *session) Reset() {

	if err := e.Parse(); err != nil {
		log.Printf("meilo: failed to parse email: %v", err)
	}

	log.Println("Sending email...")
	if err := send(e); err != nil {
		log.Printf("meilo: failed to send email: %v", err)
	}

	e.Reset()
}
