package meilo

type serverOption func(*server)

func WithPort(port string) serverOption {
	return func(s *server) {
		s.Port = port
	}
}

func WithPassword(password string) serverOption {
	return func(s *server) {
		s.Password = password
	}
}

func WithUser(user string) serverOption {
	return func(s *server) {
		s.User = user
	}
}

func WithHost(host string) serverOption {
	return func(s *server) {
		s.Host = host
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

func Only(contentTypes []string) senderOption {
	return func(s *sender) {
		s.openContentTypes = contentTypes
	}
}
