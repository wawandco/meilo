package meilo

type serverOption func(*server)

func WithPort(port string) serverOption {
	return func(s *server) {
		s.Port = port
	}
}

func WithSenderOptions(opts ...senderOption) serverOption {
	return func(s *server) {
		s.senderOpts = opts
	}
}

type senderOption func(*sender)

func WithDir(dir string) senderOption {
	return func(s *sender) {
		s.dir = dir
	}
}
