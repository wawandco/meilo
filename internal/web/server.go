package web

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/wawandco/meilo/internal/services"
	"github.com/wawandco/meilo/internal/web/templates"

	hxhttp "maragu.dev/gomponents-htmx/http"
)

type server struct {
	port         string
	httpServer   *http.Server
	emailService *services.EmailService
}

func NewServer(port string, emailService *services.EmailService) *server {
	server := &server{
		port:         port,
		emailService: emailService,
	}
	server.setupRoutes()

	return server
}

func (s *server) Start() error {
	slog.Info("Starting web server on", "port", s.port)

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("meilo: failed to start web server: %w", err)
	}

	return nil
}

func (s *server) setupRoutes() {
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /", s.list)
	mux.HandleFunc("GET /details/{id}", s.details)
	mux.HandleFunc("GET /refresh", s.refresh)
	mux.HandleFunc("DELETE /delete-all", s.clean)
	mux.HandleFunc("GET /get-latest/{email}", s.getLatestEmail)

	s.httpServer = &http.Server{
		Addr:         ":" + s.port,
		Handler:      s.loggingMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func (s *server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		wrapped := &responseWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		log.Printf("%s %s %d %v", r.Method, r.URL.Path, wrapped.statusCode, time.Since(start))
	})
}

func (s *server) list(w http.ResponseWriter, r *http.Request) {
	emails, err := s.emailService.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := templates.Layout(templates.LayoutProps{
		YieldTitle: fmt.Sprintf("Recent %v emails listed", len(emails)),
		Yield:      templates.ListEl(emails),
	}).Render(w); err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
		return
	}
}

func (s *server) refresh(w http.ResponseWriter, r *http.Request) {
	emails, err := s.emailService.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := templates.RefreshEl(fmt.Sprintf("Recent %v emails listed", len(emails)), emails).Render(w); err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
		return
	}
}

func (s *server) details(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	email, err := s.emailService.FindById(idInt64)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = templates.DetailEl(email).Render(w); err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
		return
	}
}

func (s *server) clean(w http.ResponseWriter, r *http.Request) {
	err := s.emailService.DeleteAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hxhttp.SetRedirect(w.Header(), "/")
}

func (s *server) getLatestEmail(w http.ResponseWriter, r *http.Request) {
	emailRecipient := r.PathValue("email")
	if emailRecipient == "" {
		writeResponse(w, http.StatusBadRequest, map[string]string{
			"error": "no email address provided",
		})
		return
	}
	email, err := s.emailService.GetLatestEmailByRecipient(emailRecipient)
	if err != nil {
		slog.Error("Failed to get latest email",
			"recipient", emailRecipient,
			"error", err)
		writeResponse(w, http.StatusNotFound, map[string]string{
			"error": "no email found for this recipient",
		})
		return
	}

	writeResponse(w, http.StatusOK, map[string]interface{}{
		"email": email,
	})
}
