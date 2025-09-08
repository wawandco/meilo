package smtp

import (
	"errors"
	"io"
	"log"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/wawandco/meilo/internal/web"
)

var e = email{}

type session struct {
	username, password string
	saveFn             func(e web.Email)
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

	bodies := make([]web.EmailBody, len(e.Bodies))
	for i, b := range e.Bodies {
		bodies[i] = web.EmailBody{
			ContentType: b.ContentType,
			Content:     b.Content,
		}
	}

	attachments := make([]web.EmailAttachment, len(e.Attachments))
	for i, a := range e.Attachments {
		attachments[i] = web.EmailAttachment{
			Name:        a.Name,
			Path:        a.Path,
			ContentType: a.ContentType,
			Data:        a.Data,
		}
	}

	s.saveFn(web.Email{
		Subject:     e.Subject,
		Sender:      e.From,
		Recipients:  e.To,
		CC:          e.Cc,
		BCC:         e.Bcc,
		Bodies:      bodies,
		Attachments: attachments,
		ReceivedAt:  time.Now(),
	})

	log.Println("Sending email...")
	if err := send(e); err != nil {
		log.Printf("meilo: failed to send email: %v", err)
	}

	e.Reset()
}
