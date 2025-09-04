package meilo

import "github.com/wawandco/meilo/internal/smtp"

type serverOption func(*smtp.Server)

func WithPort(port string) serverOption {
	return func(s *smtp.Server) {
		s.Port = port
	}
}

func WithDir(directory string) serverOption {
	return func(s *smtp.Server) {
		smtp.Dir = directory
	}
}

func WithWeb() serverOption {
	return func(s *smtp.Server) {
		s.EnableWeb = true
	}
}
