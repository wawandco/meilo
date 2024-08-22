package meilo

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
)

const (
	textPlain = "text/plain"
	textHtml  = "text/html"
)

type email struct {
	Subject     string
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Body        bytes.Buffer
	Bodies      []body
	Attachments []attachment
}

// Reset resets the email to its initial state.
func (e *email) Reset() {
	e.Subject = ""
	e.From = ""
	e.To = nil
	e.Cc = nil
	e.Bcc = nil
	e.Body.Reset()
	e.Bodies = nil
	e.Attachments = nil
}

// Parse parses the email and extracts the headers, bodies and attachments from it.
// It returns an error if the email could not be parsed.
func (e *email) Parse() error {
	mail, err := mail.ReadMessage(strings.NewReader(e.Body.String()))
	if err != nil {
		return fmt.Errorf("mailo: failed to parse email: %w", err)
	}

	if err := e.ParseHeaders(mail); err != nil {
		return fmt.Errorf("mailo: failed to parse headers: %w", err)
	}

	mediaType, params, err := mime.ParseMediaType(mail.Header.Get("Content-Type"))
	if err != nil {
		return fmt.Errorf("mailo: failed to parse media type: %w", err)
	}

	if strings.Contains(mediaType, "multipart") {
		if err := e.ParseMultipart(multipart.NewReader(mail.Body, params["boundary"])); err != nil {
			return fmt.Errorf("mailo: failed to parse multipart: %w", err)
		}

		return nil
	}

	fmt.Println("Parsing single part")
	e.Bodies = append(e.Bodies, body{
		ContentType: mediaType,
		Content:     e.ReadSinglePart(mail.Body),
	})

	return nil
}

// ParseHeaders parses the headers of the email and extracts the headers from it.
// It returns an error if the headers could not be parsed.
func (e *email) ParseHeaders(mail *mail.Message) error {
	e.From = mail.Header.Get("From")
	e.To = strings.Split(mail.Header.Get("To"), ",")
	e.Cc = strings.Split(mail.Header.Get("Cc"), ",")
	e.Bcc = strings.Split(mail.Header.Get("Bcc"), ",")

	//check if the subject is encoded
	subject := mail.Header.Get("Subject")
	if !strings.Contains(subject, "=?UTF-8?B?") || !strings.Contains(subject, "=?utf-8?B?") {
		e.Subject = subject
		return nil
	}

	subject = strings.ReplaceAll(subject, "=?UTF-8?B?", "")
	subject = strings.ReplaceAll(subject, "=?utf-8?B?", "")
	subject = strings.ReplaceAll(subject, "?=", "")

	decoded, err := base64.StdEncoding.DecodeString(subject)
	if err != nil {
		return fmt.Errorf("mailo: failed to decode subject: %w", err)
	}

	e.Subject = string(decoded)

	return nil
}

// ParseMultipart parses the multipart email and extracts the bodies and attachments from it.
// It returns an error if the email could not be parsed.
func (e *email) ParseMultipart(mr *multipart.Reader) error {
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("mailo: failed to read next part: %w", err)
		}

		contentType := part.Header.Get("Content-Type")
		if contentType == "" {
			continue
		}

		switch {
		case strings.Contains(contentType, textHtml):
			e.Bodies = append(e.Bodies, body{
				ContentType: textHtml,
				Content:     e.ReadPartContent(part),
			})

		case strings.Contains(contentType, textPlain):
			e.Bodies = append(e.Bodies, body{
				ContentType: textPlain,
				Content:     e.ReadPartContent(part),
			})

		default:

			if err := e.ProcessAttachments(part, contentType); err != nil {
				return fmt.Errorf("mailo: failed to process attachments: %w", err)
			}
		}

	}

	return nil
}

// ReadPartContent reads the content of the part and returns it as a string.
// It decodes the content if it is encoded in base64 or quoted-printable.
// It returns an empty string if the content could not be read.
func (e *email) ReadPartContent(part io.Reader) string {
	transferEncoding := part.(*multipart.Part).Header.Get("Content-Transfer-Encoding")
	switch {
	case strings.Contains(transferEncoding, "base64"):
		decoded, err := base64.StdEncoding.DecodeString(e.ReadSinglePart(part))
		if err != nil {
			return ""
		}
		return string(decoded)

	case strings.Contains(transferEncoding, "quoted-printable"):
		qpReader := quotedprintable.NewReader(part)
		return e.ReadSinglePart(qpReader)

	default:
		return e.ReadSinglePart(part)
	}
}

// ReadSinglePart reads the content of the part and returns it as a string.
// It returns an empty string if the content could not be read.
func (e *email) ReadSinglePart(part io.Reader) string {
	buf := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buf, part)
	if err != nil {
		log.Printf("mailo: failed to copy part: %v", err)
		return ""
	}

	return buf.String()
}

// ProcessAttachments processes the attachments of the email and appends them to the attachments slice.
// It returns an error if the attachment could not be processed.
func (e *email) ProcessAttachments(part *multipart.Part, contentType string) error {

	body, err := io.ReadAll(part)
	if err != nil {
		return fmt.Errorf("mailo: failed to read part: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return fmt.Errorf("mailo: failed to decode attachment: %w", err)
	}

	if !strings.Contains(part.FileName(), "=?UTF-8?B?") || !strings.Contains(part.FileName(), "=?utf-8?B?") {
		e.Attachments = append(e.Attachments, attachment{
			Name:        part.FileName(),
			ContentType: contentType,
			Data:        decoded,
		})

		return nil
	}

	// decoded the name of the attachment
	name := strings.ReplaceAll(part.FileName(), "=?UTF-8?B?", "")
	name = strings.ReplaceAll(name, "=?utf-8?B?", "")
	name = strings.ReplaceAll(name, "?=", "")

	decodedName, err := base64.StdEncoding.DecodeString(name)
	if err != nil {
		return fmt.Errorf("mailo: failed to decode attachment name: %w", err)
	}

	e.Attachments = append(e.Attachments, attachment{
		Name:        string(decodedName),
		ContentType: contentType,
		Data:        decoded,
	})

	return nil
}

type attachment struct {
	Name        string
	Path        string
	ContentType string
	Data        []byte
}

type body struct {
	ContentType string
	Content     string
}
