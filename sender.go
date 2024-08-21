package meilo

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"fmt"
	"html"
	"html/template"
	"log"
	"mime"
	"os"
	"path"
	"strings"

	"github.com/pkg/browser"
)

var genID = func() string {
	id := make([]byte, 16)
	_, err := rand.Read(id)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", id)
}

var (
	//go:embed html-tmpl.html
	htmlTemplate string

	//dir is the directory where the email body will be saved before opening in the browser
	dir string = os.TempDir()
)

type templConfig struct {
	Bodies      []body
	Attachments []attachment
}

func send(e email) error {
	if len(e.Bodies) == 0 {
		return fmt.Errorf("meilo: no email bodies found")
	}

	content := fmt.Sprintf(
		htmlTemplate,

		html.EscapeString(e.Subject),
		html.EscapeString(e.From),
		strings.Join(e.To, ", "),
		strings.Join(e.Cc, ", "),
		strings.Join(e.Bcc, ", "),
		html.EscapeString(e.Subject),
	)

	path, err := saveEmailBody(content, e)
	if err != nil {
		return fmt.Errorf("meilo: failed to save email body: %w", err)
	}

	if err := browser.OpenFile(path); err != nil {
		return fmt.Errorf("meilo: failed to open email in browser: %w", err)
	}

	return nil
}

func saveEmailBody(content string, email email) (string, error) {
	err := saveAttachmentFiles(email.Attachments)
	if err != nil {
		return "", fmt.Errorf("meilo: failed to save attachments: %w", err)
	}

	tmpl := template.Must(template.New("mail").Funcs(
		template.FuncMap{
			"contains": strings.Contains,
		},
	).Parse(content))

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, templConfig{
		Bodies:      email.Bodies,
		Attachments: email.Attachments,
	})

	if err != nil {
		return "", fmt.Errorf("meilo: failed to execute template: %w", err)
	}

	filePath := fmt.Sprintf("%s.html", genID())

	path := path.Join(dir, filePath)
	err = os.WriteFile(path, tpl.Bytes(), 0644)
	if err != nil {
		return "", fmt.Errorf("meilo: failed to write email body: %w", err)
	}

	return path, nil
}

func saveAttachmentFiles(attachments []attachment) error {
	for i, a := range attachments {
		if len(a.Name) > 50 {
			a.Name = a.Name[:50]
		}

		exts, err := mime.ExtensionsByType(a.ContentType)
		if err != nil {
			return fmt.Errorf("meilo: failed to get extension for attachment %s: %w", a.Name, err)
		}

		name := a.Name + genID()
		filePath := path.Join(dir, fmt.Sprintf("%s%s", name, exts[0]))

		err = os.WriteFile(filePath, a.Data, 0644)
		if err != nil {
			return fmt.Errorf("meilo: failed to write attachment %s: %w", name, err)
		}

		attachments[i].Path = filePath
	}

	return nil
}
