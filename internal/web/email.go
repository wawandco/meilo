package models

import (
	"strings"
	"time"
)

const (
	textPlain = "text/plain"
	textHtml  = "text/html"
)

type Email struct {
	ID          int64        `json:"-"`
	Subject     string       `json:"subject"`
	Sender      string       `json:"sender"`
	Recipients  []string     `json:"-"`
	CC          []string     `json:"-"`
	BCC         []string     `json:"-"`
	Body        string       `json:"emailBody"`
	Bodies      []Body       `json:"-"`
	Attachments []Attachment `json:"-"`
	ReceivedAt  time.Time    `json:"receivedAt"`
}

type Body struct {
	ContentType string
	Content     string
}

type Attachment struct {
	Name        string
	Path        string
	ContentType string
	Data        []byte
}

func (e *Email) SetHTMLBody() {
	for _, b := range e.Bodies {
		if b.ContentType == textHtml {
			e.Body = b.Content
		}
	}
}

func (e *Email) RawRecipients() string {
	return strings.Join(e.Recipients, ",")
}
