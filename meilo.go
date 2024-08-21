package meilo

import (
	"fmt"
)

// serverOption is a function that configures the server.
// It is used in the Start function.
func Start(options ...serverOption) (server, error) {
	s := server{
		Port:     "1025",
		Password: "password",
		User:     "username",
		Host:     "localhost",
	}

	for _, option := range options {
		option(&s)
	}

	go func() error {
		if err := s.run(); err != nil {
			return fmt.Errorf("meilo: failed to start server: %v", err)
		}

		return nil
	}()

	return s, nil
}
