package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Money struct {
	Value string `json:"bid"`
}

func (receiver Money) valueFormatted() string {
	return strings.Replace(`Dólar: {value}`, "value", receiver.Value, 1)
}

func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancelFunc()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", http.NoBody)
	if err != nil {
		fmt.Printf("Error on file req %s", err.Error())
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error on file req %s", err.Error())
		return
	}

	defer res.Body.Close()
	file, err := os.Create("./file")
	if err != nil {
		return
	}
	defer file.Close()

	var money Money
	err = json.NewDecoder(res.Body).Decode(&money)
	if err != nil {
		panic(err)
		return
	}
	err = json.NewEncoder(file).Encode(money.valueFormatted())
	if err != nil {
		return
	}
}
