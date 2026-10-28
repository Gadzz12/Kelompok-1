package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Inbox menampilkan halaman Inbox untuk pesan dan persetujuan (ACC) undangan kolaborasi
func (h *Handler) Inbox(w http.ResponseWriter, r *http.Request) {
	userID := h.getLoggedInUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Alihkan ke Dashboard dengan parameter tab inbox aktif
	q := r.URL.Query()
	q.Set("page", "inbox")
	http.Redirect(w, r, "/?"+q.Encode(), http.StatusSeeOther)
}

// Team menampilkan halaman daftar anggota tim dan kolaborator workspace
func (h *Handler) Team(w http.ResponseWriter, r *http.Request) {
	userID := h.getLoggedInUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	q := r.URL.Query()
	q.Set("page", "team")
	http.Redirect(w, r, "/?"+q.Encode(), http.StatusSeeOther)
}

// Documents menampilkan halaman berkas dan dokumen referensi tugas
func (h *Handler) Documents(w http.ResponseWriter, r *http.Request) {
	userID := h.getLoggedInUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	q := r.URL.Query()
	q.Set("page", "documents")
	http.Redirect(w, r, "/?"+q.Encode(), http.StatusSeeOther)
}

// Calendar menampilkan jadwal deadline tugas berdasarkan tanggal
func (h *Handler) Calendar(w http.ResponseWriter, r *http.Request) {
	userID := h.getLoggedInUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	q := r.URL.Query()
	q.Set("page", "calendar")
	http.Redirect(w, r, "/?"+q.Encode(), http.StatusSeeOther)
}

// AttachDocument menangani penambahan atau pengunggahan dokumen ke tugas tertentu
func (h *Handler) AttachDocument(w http.ResponseWriter, r *http.Request) {
	userID := h.getLoggedInUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Dukung multipart form jika ada berkas yang diunggah
	_ = r.ParseMultipartForm(10 << 20)

	todoIDStr := strings.TrimSpace(r.FormValue("todo_id"))
	docName := strings.TrimSpace(r.FormValue("document"))
	fromPage := strings.TrimSpace(r.FormValue("from_page"))

	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil || todoID <= 0 {
		h.redirectWithError(w, r, fromPage, todoID, "Please select a valid task to attach document.")
		return
	}

	file, header, fileErr := r.FormFile("document_file")
	if fileErr == nil && header != nil {
		defer file.Close()
		docName = header.Filename
	}

	if docName == "" {
		h.redirectWithError(w, r, fromPage, todoID, "Document name or file cannot be empty.")
		return
	}

	if err := h.repo.AttachDocument(todoID, userID, docName); err != nil {
		h.redirectWithError(w, r, fromPage, todoID, err.Error())
		return
	}

	msg := fmt.Sprintf("Document '%s' successfully attached to task.", docName)
	if fromPage == "detail" {
		http.Redirect(w, r, fmt.Sprintf("/todo/detail?id=%d&msg=%s", todoID, url.QueryEscape(msg)), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?page=documents&msg="+url.QueryEscape(msg), http.StatusSeeOther)
}

// SendInvitation mengirimkan undangan kolaborasi tugas ke pengguna lain
func (h *Handler) SendInvitation(w http.ResponseWriter, r *http.Request) {
	userID := h.getLoggedInUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	todoIDStr := strings.TrimSpace(r.FormValue("todo_id"))
	invitee := strings.TrimSpace(r.FormValue("invitee"))
	fromPage := strings.TrimSpace(r.FormValue("from_page"))

	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil || todoID <= 0 {
		h.redirectWithError(w, r, fromPage, todoID, "Please select a valid task for collaboration.")
		return
	}

	if invitee == "" {
		h.redirectWithError(w, r, fromPage, todoID, "Collaborator email or username is required.")
		return
	}

	_, err = h.repo.InviteCollaborator(todoID, userID, invitee)
	if err != nil {
		h.redirectWithError(w, r, fromPage, todoID, err.Error())
		return
	}

	// Berhasil mengirim undangan
	msg := fmt.Sprintf("Collaboration invitation successfully sent to %s. Awaiting approval (ACC) via inbox.", invitee)
	if fromPage == "detail" {
		http.Redirect(w, r, fmt.Sprintf("/todo/detail?id=%d&msg=%s", todoID, url.QueryEscape(msg)), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?page=team&msg="+url.QueryEscape(msg), http.StatusSeeOther)
}

// RespondInvitation menangani konfirmasi persetujuan (ACC) atau penolakan undangan kolaborasi
func (h *Handler) RespondInvitation(w http.ResponseWriter, r *http.Request) {
	userID := h.getLoggedInUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	inviteID, err := strconv.Atoi(r.FormValue("invite_id"))
	if err != nil || inviteID <= 0 {
		http.Redirect(w, r, "/?page=inbox&error="+url.QueryEscape("Invalid invitation ID"), http.StatusSeeOther)
		return
	}

	action := strings.TrimSpace(strings.ToLower(r.FormValue("action")))
	accept := action == "accept"

	user, exists := h.repo.GetUserByID(userID)
	if !exists {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	err = h.repo.RespondInvitation(inviteID, userID, user.Email, accept)
	if err != nil {
		log.Println("Gagal merespons undangan kolaborasi:", err)
		http.Redirect(w, r, "/?page=inbox&error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	var successMsg string
	if accept {
		successMsg = "Collaboration invitation accepted (ACC). Task is now available in your workspace!"
	} else {
		successMsg = "Collaboration invitation rejected."
	}

	http.Redirect(w, r, "/?page=inbox&msg="+url.QueryEscape(successMsg), http.StatusSeeOther)
}

// MarkAllInboxRead menandai seluruh pesan inbox pengguna sebagai sudah dibaca
func (h *Handler) MarkAllInboxRead(w http.ResponseWriter, r *http.Request) {
	userID := h.getLoggedInUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, exists := h.repo.GetUserByID(userID)
	if !exists {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	if err := h.repo.MarkAllInboxAsRead(user.Email); err != nil {
		log.Println("Gagal menandai inbox terbaca:", err)
	}

	http.Redirect(w, r, "/?page=inbox&msg="+url.QueryEscape("All inbox notifications have been marked as read."), http.StatusSeeOther)
}

func (h *Handler) redirectWithError(w http.ResponseWriter, r *http.Request, fromPage string, todoID int, errMsg string) {
	if fromPage == "detail" && todoID > 0 {
		http.Redirect(w, r, fmt.Sprintf("/todo/detail?id=%d&error=%s", todoID, url.QueryEscape(errMsg)), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?page=team&error="+url.QueryEscape(errMsg), http.StatusSeeOther)
}


// AddFriend menangani penambahan teman ke dalam tim
func (h *Handler) AddFriend(w http.ResponseWriter, r *http.Request) {
	userID := h.getLoggedInUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	identifier := strings.TrimSpace(r.FormValue("identifier"))
	if identifier == "" {
		http.Redirect(w, r, "/?page=team&error="+url.QueryEscape("Please enter an email or username."), http.StatusSeeOther)
		return
	}

	addedUser, err := h.repo.AddFriend(userID, identifier)
	if err != nil {
		http.Redirect(w, r, "/?page=team&error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	msg := fmt.Sprintf("✓ Successfully added %s (@%s) to your Team!", addedUser.Nama, addedUser.Username)
	http.Redirect(w, r, "/?page=team&msg="+url.QueryEscape(msg), http.StatusSeeOther)
}

// RemoveFriend menangani penghapusan teman dari tim
func (h *Handler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	userID := h.getLoggedInUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	friendID, err := strconv.Atoi(r.FormValue("friend_id"))
	if err != nil || friendID <= 0 {
		http.Redirect(w, r, "/?page=team&error="+url.QueryEscape("Invalid friend ID."), http.StatusSeeOther)
		return
	}

	if err := h.repo.RemoveFriend(userID, friendID); err != nil {
		http.Redirect(w, r, "/?page=team&error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	msg := "Friend removed from your Team."
	http.Redirect(w, r, "/?page=team&msg="+url.QueryEscape(msg), http.StatusSeeOther)
}
