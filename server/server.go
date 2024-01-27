package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"time"
)

const BASE_URL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type ResponseApi struct {
	USDBRL Usdbrl
}

type Usdbrl struct {
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

type ResponseClient struct {
	Bid string `json:"bid"`
}

func init() {
	db, err := openConnectionDB()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS usdbrl (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				code TEXT,
				codein TEXT,
				name TEXT,
				high TEXT,
				low TEXT,
				varBid TEXT,
				pctChange TEXT,
				bid TEXT,
				ask TEXT,
				timestamp TEXT,
				create_date TEXT
			);
`)
	if err != nil {
		panic(err)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		money, err := getMoney(saveOnDb)
		if err != nil {
			fmt.Printf("Error get money %s \r\n", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(ResponseClient{Bid: money.USDBRL.Bid})
		if err != nil {
			fmt.Printf("Error on convert money %s \r\n", err.Error())
			return
		}
	})
	http.ListenAndServe(":8080", mux)
}

func getMoney(callback func(responseApi *ResponseApi)) (*ResponseApi, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancelFunc()
	req, err := http.NewRequestWithContext(ctx, "GET", BASE_URL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("Error create request %s \r\n", err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error get money %s \r\n", err.Error())
	}

	defer res.Body.Close()

	var response ResponseApi
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("Error decode money %s \r\n", err.Error())
	}

	go callback(&response)
	return &response, nil
}

func saveOnDb(responseApi *ResponseApi) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancelFunc()

	db, err := openConnectionDB()
	if err != nil {
		return
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, `
	INSERT INTO usdbrl (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, responseApi.USDBRL.Code,
		responseApi.USDBRL.Codein,
		responseApi.USDBRL.Name,
		responseApi.USDBRL.High,
		responseApi.USDBRL.Low,
		responseApi.USDBRL.VarBid,
		responseApi.USDBRL.PctChange,
		responseApi.USDBRL.Bid,
		responseApi.USDBRL.Ask,
		responseApi.USDBRL.Timestamp,
		responseApi.USDBRL.CreateDate)

	if err != nil {
		fmt.Printf("Error on save %s", err.Error())
		return
	}
}

func openConnectionDB() (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", "./sqlite-db.db")
	if err != nil {
		return nil, err
	}

	return conn, nil
}
