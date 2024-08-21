package meilo

type serverOption func(*server)

func WithPort(port string) serverOption {
	return func(s *server) {
		s.Port = port
	}
}

func WithDir(directory string) serverOption {
	return func(s *server) {
		dir = directory
	}
}
