package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	TEN_MILLISECONDS         = 10 * time.Millisecond
	TWO_HUNDRED_MILLISECONDS = 200 * time.Millisecond
)

type Quotation struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	db, err := sql.Open("sqlite3", "./db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		FindQuotationHandler(w, r, db)
	})
	http.ListenAndServe(":8080", nil)

}

func SetupBD() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./db")
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(
		"CREATE TABLE IF NOT EXISTS" +
			" quotations (id INTEGER PRIMARY KEY, code TEXT, code_in TEXT," +
			" name TEXT, high TEXT, low TEXT," +
			" var_bid TEXT, pct_change TEXT, bid TEXT," +
			" ask TEXT, timestamp TEXT, created_at TEXT)")

	if err != nil {
		return nil, err
	}

	if _, err := stmt.Exec(); err != nil {
		return nil, err
	}

	return db, nil
}

func FindQuotationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	quotation, err := FindQuotation(r.Context(), db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), TEN_MILLISECONDS)
	defer cancel()

	if err := insertQuotation(ctx, db, quotation); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Tempo insuficiente ao inserir no banco de dados")
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quotation)
}

func FindQuotation(ctx context.Context, db *sql.DB) (*Quotation, error) {
	data, err := RequestQuotation(ctx)
	if err != nil {
		return nil, err
	}

	value, err := json.Marshal((*data)["USDBRL"])
	if err != nil {
		return nil, err
	}

	var quotation Quotation
	if err := json.Unmarshal(value, &quotation); err != nil {
		return nil, err
	}

	return &quotation, nil
}

func RequestQuotation(ctx context.Context) (*map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, TWO_HUNDRED_MILLISECONDS)
	defer cancel()

	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Tempo insuficiente ao buscar na API")
		}
		return nil, err
	}
	defer res.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func insertQuotation(ctx context.Context, db *sql.DB, quotation *Quotation) error {
	stmt, err := db.Prepare(
		" INSERT INTO" +
			" quotations(code, code_in, name, high, low," +
			" var_bid, pct_change, bid, ask, timestamp, created_at)" +
			" VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		quotation.Code,
		quotation.Codein,
		quotation.Name,
		quotation.High,
		quotation.Low,
		quotation.VarBid,
		quotation.PctChange,
		quotation.Bid,
		quotation.Ask,
		quotation.Timestamp,
		quotation.CreateDate,
	)
	if err != nil {
		return err
	}

	return nil
}
