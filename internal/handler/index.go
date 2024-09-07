package handler

import (
	"database/sql"
	"fmt"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	db "shareem/database/sqlc"
	"shareem/internal/share"
	"strings"

	goaway "github.com/TwiN/go-away"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Server struct {
	logger *slog.Logger
	tmpl   *template.Template
	store  db.Querier
}

func NewServer(logger *slog.Logger, tmpl *template.Template, dbPool *pgxpool.Pool) *Server {
	store := db.New(dbPool)
	return &Server{
		logger: logger,
		tmpl:   tmpl,
		store:  store,
	}
}

// used for the template index page to be filled
type indexPageData struct {
	Shares []db.Share
	Total  int64
}

type errorPage struct {
	ErrorMsg string
}

// serves the index page : shows all the shares
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	// TODO: add limit to the fetched data
	shares, err := s.store.ListShares(r.Context())

	if err != nil {
		s.logger.Error("could not fetch shares", slog.Any("Error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	total, err := s.store.CountShares(r.Context())
	if err != nil {
		s.logger.Error("could not fetch shares", slog.Any("Error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	s.tmpl.ExecuteTemplate(w, "index.html", indexPageData{
		Shares: shares,
		Total:  total,
	})
}

// serves the add share page: adds a new share
func (s *Server) Insert(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.logger.Error("could not parse form", slog.Any("Error", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, ok := r.Form["url"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	urlStr := strings.Join(url, "")

	if strings.TrimSpace(urlStr) == "" {
		w.WriteHeader(http.StatusBadRequest)
		s.tmpl.ExecuteTemplate(w, "error.html", errorPage{
			ErrorMsg: "URL cannot be empty",
		})
		return
	}

	note, ok := r.Form["note"]
	noteStr := strings.Join(note, "")

	// Convert noteStr to sql.NullString
	var dbNote sql.NullString
	if strings.TrimSpace(noteStr) == "" {
		dbNote = sql.NullString{String: "", Valid: false}
	} else {
		dbNote = sql.NullString{String: noteStr, Valid: true}
	}

	splits := strings.Split(r.RemoteAddr, ":")
	ipStr := strings.Trim(strings.Join(splits[:len(splits)-1], ":"), "[]")
	ip := net.ParseIP(ipStr)

	// Convert IP to pgtype.Inet
	var dbIP pgtype.Inet
	err := dbIP.Set(ip)
	if err != nil {
		s.logger.Error("could not parse IP", slog.Any("Error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if goaway.IsProfane(noteStr) {
		w.WriteHeader(http.StatusBadRequest)
		s.tmpl.ExecuteTemplate(w, "error.html", errorPage{
			ErrorMsg: fmt.Sprintf(
				"Please don't use profanity. Your IP has been tracked %s",
				ipStr,
			),
		})
		return
	}

	share, err := share.NewShare(urlStr, noteStr, ip)
	if err != nil {
		s.logger.Error("could not create share", slog.Any("Error", err))

		// use the error page Template
		w.WriteHeader(http.StatusBadRequest)
		s.tmpl.ExecuteTemplate(w, "error.html", errorPage{
			ErrorMsg: err.Error(),
		})
		return
	}

	dbShare := db.CreateShareParams{
		ID:        share.ID,
		Url:       share.URL,
		Note:      dbNote,
		CreatedAt: share.CreatedAt,
		Ip:        dbIP,
		UpdatedAt: share.CreatedAt,
	}

	if _, err := s.store.CreateShare(r.Context(), dbShare); err != nil {
		s.logger.Error("could not insert share", slog.Any("Error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
