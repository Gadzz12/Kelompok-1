package handler

import (
	"encoding/base64"
	"html/template"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"todo-app-test/internal/repository"
)

// embeddedLogoBase64 menyimpan aset logo resmi Flowspace secara mandiri dalam program
// sehingga tidak memerlukan file fisik logo.png pada disk.
const embeddedLogoBase64 = "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"

var embeddedLogoBytes []byte

func init() {
	var err error
	embeddedLogoBytes, err = base64.StdEncoding.DecodeString(embeddedLogoBase64)
	if err != nil {
		panic(err)
	}
}

// Handler mengelola HTTP request untuk aplikasi Todo
type Handler struct {
	repo repository.Repository
	tmpl *template.Template
}

// New membuat instance Handler baru
func New(repo repository.Repository, tmpl *template.Template) *Handler {
	return &Handler{
		repo: repo,
		tmpl: tmpl,
	}
}

// RegisterRoutes mendaftarkan seluruh endpoint aplikasi ke router
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/register", h.Register)
	mux.HandleFunc("/logout", h.Logout)
	mux.HandleFunc("/", h.Dashboard)
	mux.HandleFunc("/tasks", h.Dashboard)
	mux.HandleFunc("/settings", h.Settings)
	mux.HandleFunc("/user/preferences", h.UpdateUserPreferences)
	mux.HandleFunc("/todo/create", h.CreateTodo)
	mux.HandleFunc("/todo/detail", h.DetailTodo)
	mux.HandleFunc("/todo/edit", h.EditTodo)
	mux.HandleFunc("/todo/delete", h.DeleteTodo)

	// Halaman & Fitur Workspace (Inbox, Team, Document, Calendar)
	mux.HandleFunc("/inbox", h.Inbox)
	mux.HandleFunc("/team", h.Team)
	mux.HandleFunc("/documents", h.Documents)
	mux.HandleFunc("/calendar", h.Calendar)
	mux.HandleFunc("/todo/attach-document", h.AttachDocument)
	mux.HandleFunc("/invitation/send", h.SendInvitation)
	mux.HandleFunc("/invitation/respond", h.RespondInvitation)
	mux.HandleFunc("/inbox/read-all", h.MarkAllInboxRead)
	mux.HandleFunc("/friend/add", h.AddFriend)
	mux.HandleFunc("/friend/remove", h.RemoveFriend)

	// File statis (CSS, JS) untuk dashboard dan halaman lainnya
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	serveLogo := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Write(embeddedLogoBytes)
	}
	mux.HandleFunc("/static/logo.png", serveLogo)
	mux.HandleFunc("/logo.png", serveLogo)
}

// getLoggedInUserID mengambil user ID dari cookie sesi
func (h *Handler) getLoggedInUserID(r *http.Request) int {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return 0
	}

	id, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return 0
	}

	return id
}

// emailIsValid memeriksa kevalidan format alamat email
func emailIsValid(email string) bool {
	email = strings.TrimSpace(email)
	address, err := mail.ParseAddress(email)
	return err == nil && address.Address == email
}
