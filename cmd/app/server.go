package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/reyhardy/exs-prescription-app/internal/db"
	"github.com/reyhardy/exs-prescription-app/internal/handler"
	"github.com/reyhardy/exs-prescription-app/internal/handler/middleware"
	"github.com/reyhardy/exs-prescription-app/internal/repo"
	"github.com/reyhardy/exs-prescription-app/internal/service"
	"github.com/reyhardy/exs-prescription-app/internal/template"
)

type Repos struct {
	account   *repo.Account
	session   *repo.Session
	patient   *repo.Patient
	complaint *repo.Complaint
}

func NewRepos(db *db.DB) Repos {
	return Repos{
		account:   repo.NewAccount(db),
		session:   repo.NewSession(db),
		patient:   repo.NewPatient(db),
		complaint: repo.NewComplaint(db),
	}
}

type Services struct {
	account *service.User
	sm      *service.SessionManager
	patient *service.Patient
}

func NewServices(r Repos) Services {
	return Services{
		account: service.NewUser(r.account),
		sm:      service.NewSession(r.session, 20*time.Second),
		patient: service.NewPatient(r.patient, r.complaint),
	}
}

type Handlers struct {
	mw middleware.Middleware
	*handler.User
	*handler.Auth
	*handler.Patient
}

func NewHandler(s Services) (Handlers, error) {
	tm, err := template.NewManager()
	if err != nil {
		return Handlers{}, err
	}

	return Handlers{
		mw:      middleware.New(s.account, s.sm),
		User:    handler.NewUser(s.account, tm),
		Auth:    handler.NewAuth(s.account, s.sm, tm),
		Patient: handler.NewPatient(s.patient, tm),
	}, nil
}

func setupRoutes(h Handlers) chi.Router {
	mux := chi.NewMux()
	mux.Use(chimw.Recoverer)
	mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("public/"))))
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
	mux.HandleFunc("GET /login", h.LoginPage)
	mux.HandleFunc("POST /login", h.UserLogin)
	mux.Group(func(r chi.Router) {
		r.Use(h.mw.Authentication)
		r.HandleFunc("GET /profile", h.UserProfilePage)
		r.HandleFunc("POST /logout", h.UserLogout)
	})
	mux.Group(func(r chi.Router) {
		r.Use(h.mw.Authentication)
		r.Use(h.mw.RBAC("therapist"))
		r.Route("/patients", func(r chi.Router) {
			r.HandleFunc("GET /", h.ListAllPatients)
			r.HandleFunc("GET /{ptId}", h.ViewPatientInfo)
		})
	})
	return mux
}

func run(db *db.DB) error {
	r := NewRepos(db)
	s := NewServices(r)
	h, err := NewHandler(s)
	if err != nil {
		return err
	}
	mux := setupRoutes(h)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		ErrorLog: slog.NewLogLogger(
			slog.NewJSONHandler(os.Stdout, nil),
			slog.LevelInfo,
		),
	}
	log.Println("Listening to port 8080")
	return srv.ListenAndServe()
}
