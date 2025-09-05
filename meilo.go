package meilo

import (
	"fmt"
	"log"

	"github.com/wawandco/meilo/internal/smtp"
	"github.com/wawandco/meilo/internal/storage"
	"github.com/wawandco/meilo/internal/web"
)

// Start initializes an SMTP server with the provided configuration options.
// Each serverOption applies specific settings during server creation.
func Start(options ...serverOption) (smtp.Server, error) {
	emailStorage := storage.NewMemoryEmailStore()

	s := smtp.Server{
		Port:     "1025",
		Password: "password",
		User:     "username",
		Host:     "localhost",
		WebPort:  "",
		SaveFn:   emailStorage.Add,
	}

	for _, option := range options {
		option(&s)
	}

	go func() {
		err := func() error {
			if err := s.Run(); err != nil {
				return fmt.Errorf("meilo: failed to start server: %v", err)
			}

			return nil
		}()
		if err != nil {
			log.Printf("meilo: failed to start server: %v", err)
		}
	}()

	if s.WebPort != "" {
		webServer := web.NewServer(s.WebPort, emailStorage)
		go func() {
			err := webServer.Start()
			if err != nil {
				log.Printf("meilo: failed to start server: %v", err)
				return
			}
		}()
	}

	return s, nil
}
