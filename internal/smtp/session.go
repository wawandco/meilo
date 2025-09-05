package smtp

import (
	"errors"
	"io"
	"log"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/wawandco/meilo/internal/models"
)

var e = email{}

type session struct {
	username, password string
	saveFn             func(e models.Email)
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

	bodies := make([]models.Body, len(e.Bodies))
	for i, b := range e.Bodies {
		bodies[i] = models.Body{
			ContentType: b.ContentType,
			Content:     b.Content,
		}
	}

	attachments := make([]models.Attachment, len(e.Attachments))
	for i, a := range e.Attachments {
		attachments[i] = models.Attachment{
			Name:        a.Name,
			Path:        a.Path,
			ContentType: a.ContentType,
			Data:        a.Data,
		}
	}

	s.saveFn(models.Email{
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
