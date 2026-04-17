package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type Expense struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
}

// type DeleteResponse struct {
// 	Success bool `json:"success"`
// }

// type SummaryResponse struct {
// 	Month int     `json:"month,omitempty"`
// 	Total float64 `json:"total"`
// }

// var expenses []Expense
// var nextID = 1

// POST
func addExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		json.NewEncoder(w).Encode(ErrorResponse{"dine POST method isleyar", "400"})
		return
	}

	var req Expense

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{err.Error(), "400"})
		return
	}
	id_string := strconv.Itoa(req.ID)
	amount_string := strconv.FormatFloat(req.Amount, 'f', -1, 64)

	// expenses = append(expenses, exp)

	_, err = db.Exec(context.Background(), "insert into expenses (id, date, description, amount) values ( '"+id_string+"', '"+req.Description+"', '"+amount_string+"', '"+req.Date+"');")
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{err.Error(), "400"})
		return
	}

	json.NewEncoder(w).Encode(true)
	// if err != nil {
	// 	fmt.Printf("Error encoding response: %v\n", err)
	// }
}

// GET
// func getExpenses(w http.ResponseWriter, r *http.Request) {
// 	rows, err := db.Query(context.Background(), "select id, date, description, amount from expenses")
// 	if err != nil {
// 		fmt.Fprintln(w, "YALNYSLYK: ", err.Error())
// 	}
// 	var list []Expense
// 	for rows.Next() {
// 		var res Expense
// 		err = rows.Scan(&res.ID, &res.Description, &res.Amount, &res.Date)
// 		if err != nil {
// 			fmt.Fprintln(w, "YALNYSLYK: ", err.Error())
// 		}
// 		list = append(list, res)
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(expenses)
// }

// // DELETE
// func deleteExpense(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")
// 	id, _ := strconv.Atoi(idStr)

// 	for i, e := range expenses {
// 		if e.ID == id {
// 			expenses = append(expenses[:i], expenses[i+1:]...)
// 			json.NewEncoder(w).Encode(DeleteResponse{Success: true})
// 			return
// 		}
// 	}

// 	json.NewEncoder(w).Encode(DeleteResponse{Success: false})
// }

// // UPDATE
// func updateExpense(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")
// 	id, _ := strconv.Atoi(idStr)

// 	var req Expense
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		fmt.Fprintln(w, "YALNYSLYK: ", err.Error())
// 		return
// 	}

// 	for i, e := range expenses {
// 		if e.ID == id {
// 			expenses[i].Description = req.Description
// 			expenses[i].Amount = req.Amount

// 			json.NewEncoder(w).Encode(expenses[i])
// 			return
// 		}
// 	}

// }

// func summary(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	monthStr := r.URL.Query().Get("month")
// 	total := 0.0

// 	if monthStr == "" {
// 		for _, e := range expenses {
// 			total += e.Amount
// 		}

// 		json.NewEncoder(w).Encode(SummaryResponse{
// 			Total: total,
// 		})
// 		return
// 	}

// 	month, _ := strconv.Atoi(monthStr)

// 	for _, e := range expenses {
// 		t, _ := time.Parse("2006-01-02", e.Date)
// 		if int(t.Month()) == month && t.Year() == time.Now().Year() {
// 			total += e.Amount
// 		}
// 	}

// 	json.NewEncoder(w).Encode(SummaryResponse{
// 		Month: month,
// 		Total: total,
// 	})
// }

// func method(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		addExpense(w, r)
// 	} else if r.Method == "GET" {
// 		getExpenses(w, r)
// 	} else if r.Method == "DELETE" {
// 		deleteExpense(w, r)
// 	} else if r.Method == "PUT" {
// 		updateExpense(w, r)
// 	}
// }

func main() {
	// http.HandleFunc("/expenses/summary", summary)

	http.HandleFunc("/expenses", addExpense)

	http.ListenAndServe(":8080", nil)
}

func connectDB(config string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}
