package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os/exec"
)

const (
	AdminPassword = "admin1234"
	SecretKey     = "sk-live-9f8e7d6c5b4a3210"
	DBConnString  = "postgres://admin:password123@prod-db/users"
)

// Login authenticates a user — hardcoded credentials + SQL injection
func Login(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// CWE-798: hardcoded fallback admin credential
	if password == AdminPassword {
		w.Write([]byte("admin access granted"))
		return
	}

	// CWE-89: SQL injection — user input concatenated directly into query
	query := fmt.Sprintf("SELECT id FROM users WHERE username='%s' AND password='%s'", username, password)
	row := db.QueryRow(query)
	var id int
	if err := row.Scan(&id); err != nil {
		http.Error(w, "unauthorized", 401)
		return
	}
	w.Write([]byte(fmt.Sprintf("welcome user %d", id)))
}

// RunDiag runs a server diagnostic using user-supplied input
func RunDiag(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")

	// CWE-78: OS command injection — attacker can pass host=127.0.0.1;cat /etc/passwd
	out, err := exec.Command("bash", "-c", "ping -c 1 "+host).Output()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out)
}

// GetFile serves user files with no path sanitisation
func GetFile(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")

	// CWE-22: path traversal — attacker can read arbitrary files with ../../etc/passwd
	http.ServeFile(w, r, "/var/uploads/"+name)
}
