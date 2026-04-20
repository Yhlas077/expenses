package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type Expense struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
}

type SummaryResponse struct {
	Month int `json:"month,omitempty"`
	Total int `json:"total"`
}

// CONNECT DB
func connectDB(config string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB error: %v\n", err)
		os.Exit(1)
	}
	return conn
}

// ADD EXPENSE (POST)
func addExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Date        string `json:"date"`
		Description string `json:"description"`
		Amount      int    `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{err.Error(), "400"})
		return
	}

	// parse date
	t, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{"invalid date format", "400"})
		return
	}

	_, err = db.Exec(context.Background(),
		"INSERT INTO expenses(date, description, amount) VALUES ($1, $2, $3)",
		t, req.Description, req.Amount,
	)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorRespssonse{err.Error(), "400"})
		return
	}

	json.NewEncoder(w).Encode(true)
}

// GET EXPENSES
func getExpenses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query(context.Background(),
		"SELECT id, date, description, amount FROM expenses ORDER BY id DESC")
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{err.Error(), "500"})
		return
	}
	defer rows.Close()

	var list []Expense

	for rows.Next() {
		var e Expense

		err = rows.Scan(&e.ID, &e.Date, &e.Description, &e.Amount)
		if err != nil {
			json.NewEncoder(w).Encode(ErrorResponse{err.Error(), "500"})
			return
		}

		list = append(list, e)
	}

	json.NewEncoder(w).Encode(list)
}

// DELETE EXPENSE
func deleteExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(false)
		return
	}

	res, err := db.Exec(context.Background(),
		"DELETE FROM expenses WHERE id=$1", id)

	if err != nil || res.RowsAffected() == 0 {
		json.NewEncoder(w).Encode(false)
		return
	}

	json.NewEncoder(w).Encode(true)
}

// UPDATE EXPENSE
func updateExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}

	var req struct {
		Description string `json:"description"`
		Amount      int    `json:"amount"`
		Date        string `json:"date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{err.Error(), "400"})
		return
	}

	_, err = db.Exec(context.Background(),
		"UPDATE expenses SET description=$1, amount=$2 WHERE id=$3",
		req.Description, req.Amount, id,
	)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{err.Error(), "400"})
		return
	}

	json.NewEncoder(w).Encode(true)
}

// SUMMARY
func summary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	monthStr := r.URL.Query().Get("month")
	total := 0

	rows, err := db.Query(context.Background(),
		"SELECT date, amount FROM expenses")
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{err.Error(), "500"})
		return
	}
	defer rows.Close()

	var month int
	if monthStr != "" {
		month, _ = strconv.Atoi(monthStr)
	}

	nowYear := time.Now().Year()

	for rows.Next() {
		var date time.Time
		var amount int

		rows.Scan(&date, &amount)

		if monthStr == "" || (int(date.Month()) == month && date.Year() == nowYear) {
			total += amount
		}
	}

	json.NewEncoder(w).Encode(SummaryResponse{
		Month: month,
		Total: total,
	})
}

// ROUTER
func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getExpenses(w, r)
	case "POST":
		addExpense(w, r)
	case "DELETE":
		deleteExpense(w, r)
	case "PUT":
		updateExpense(w, r)
	}
}

// MAIN
func main() {
	db = connectDB("postgres://yhlas1:123456@localhost:5432/postgres")
	defer db.Close(context.Background())

	http.HandleFunc("/expenses", handler)
	http.HandleFunc("/expenses/summary", summary)

	http.ListenAndServe(":8080", nil)
}
