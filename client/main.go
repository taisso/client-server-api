package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	THREE_HUNDRED_MILLISECONDS = 300 * time.Millisecond
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
	ctx, cancel := context.WithTimeout(context.Background(), THREE_HUNDRED_MILLISECONDS)
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	defer cancel()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Tempo insuficiente ao buscar no server")
			return
		}
		panic(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var quotation Quotation
	if err := json.Unmarshal(body, &quotation); err != nil {
		panic(err)
	}

	if err := saveFile(quotation.Bid); err != nil {
		panic(err)
	}

	fmt.Println(string(body))

}

func saveFile(bid string) error {
	f, err := os.Create("cotacao.txt")
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("Dólar: %s", bid))
	if err != nil {
		return err
	}

	return nil
}
