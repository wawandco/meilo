package meilo

import (
	"fmt"
	"log"

	"github.com/wawandco/meilo/internal/database"
	"github.com/wawandco/meilo/internal/services"
	"github.com/wawandco/meilo/internal/smtp"
	"github.com/wawandco/meilo/internal/web"
)

// Start initializes an SMTP server with the provided configuration options.
// Each serverOption applies specific settings during server creation.
func Start(options ...serverOption) (smtp.Server, error) {
	// initialize database
	dbConnection := database.Initialize()

	// run migrations
	err := database.RunMigrations(dbConnection)
	if err != nil {
		log.Fatal(err)
	}

	emailService := services.NewEmailService(dbConnection)

	s := smtp.Server{
		Port:     "1025",
		Password: "password",
		User:     "username",
		Host:     "localhost",
		SaveFn:   emailService.Save,
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

	if s.EnableWeb {
		webServer := web.NewServer("8080", emailService)
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
