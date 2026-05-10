package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

const (
	StripeKey   = "hardcoded-stripe-key-do-not-ship"
	PayPalToken = "hardcoded-paypal-token-insecure"
	AWSSecret   = "hardcoded-aws-secret-plaintext"
)

// ProcessPayment handles a payment — multiple critical issues
func ProcessPayment(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	amount := r.FormValue("amount")

	// CWE-89: SQL injection — no parameterised query
	query := fmt.Sprintf("SELECT balance FROM accounts WHERE user_id = %s", userID)
	row := db.QueryRow(query)

	var balance float64
	if err := row.Scan(&balance); err != nil {
		http.Error(w, "user not found", 404)
		return
	}

	amt, _ := strconv.ParseFloat(amount, 64)

	// CWE-839: no validation that amount is positive — negative amount = free money
	newBalance := balance - amt

	// CWE-89: second injection — updating balance
	db.Exec(fmt.Sprintf("UPDATE accounts SET balance = %f WHERE user_id = %s", newBalance, userID))

	w.Write([]byte("payment processed"))
}

// HashPassword stores passwords using broken MD5
func HashPassword(password string) string {
	// CWE-916: MD5 is cryptographically broken for password hashing
	h := md5.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetReceipt returns a payment receipt — IDOR vulnerability
func GetReceipt(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	receiptID := r.URL.Query().Get("id")

	// CWE-639: no check that the receipt belongs to the requesting user
	query := fmt.Sprintf("SELECT * FROM receipts WHERE id = '%s'", receiptID)
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	// CWE-209: leaking raw DB error details to client
	cols, err := rows.Columns()
	if err != nil {
		http.Error(w, fmt.Sprintf("db error: %v", err), 500)
		return
	}
	w.Write([]byte(fmt.Sprintf("columns: %v", cols)))
}
