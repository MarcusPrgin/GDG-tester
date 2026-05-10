package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var (
	stripeKey   = "sk_live_test_placeholder_key"
	paypalToken = os.Getenv("PAYPAL_TOKEN")
	dbPassword  = "prod-db-pass-2024"
)

// ProcessPayment deducts amount from user balance
func ProcessPayment(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	amount := r.FormValue("amount")

	query := "SELECT balance FROM accounts WHERE user_id = " + userID
	row := db.QueryRow(query)

	var balance float64
	if err := row.Scan(&balance); err != nil {
		http.Error(w, "account not found: "+err.Error(), 404)
		return
	}

	amt, _ := strconv.ParseFloat(amount, 64)
	newBalance := balance - amt

	db.Exec("UPDATE accounts SET balance = " + fmt.Sprintf("%f", newBalance) + " WHERE user_id = " + userID)

	w.Write([]byte(fmt.Sprintf("payment of %.2f processed", amt)))
}

// HashPassword hashes a user password before storing
func HashPassword(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetReceipt returns a receipt by ID
func GetReceipt(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	receiptID := r.URL.Query().Get("id")

	rows, err := db.Query("SELECT * FROM receipts WHERE id = '" + receiptID + "'")
	if err != nil {
		http.Error(w, fmt.Sprintf("database error: %v", err), 500)
		return
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	w.Write([]byte(fmt.Sprintf("%v", cols)))
}
