package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CurrencyInfo struct {
	Date  string  `json:"date"`
	Order string  `json:"order"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type AlfaExchange struct {
	Request struct {
			Filter struct {
				       Segment string `json:"segment"`
				       Text    string `json:"text"`
			       } `json:"filter"`
			Limit   string `json:"limit"`
			Offset  string `json:"offset"`
			Order   string `json:"order"`
			Server  string `json:"server"`
			Service string `json:"service"`
			Version string `json:"version"`
		} `json:"request"`
	Response struct {
			Data struct {
				     Chf []CurrencyInfo `json:"chf"`
				     Eur []CurrencyInfo `json:"eur"`
				     Gbp []CurrencyInfo `json:"gbp"`
				     Usd []CurrencyInfo `json:"usd"`
			     } `json:"data"`
			Status string `json:"status"`
		} `json:"response"`
}


func main() {
	resp, err := http.Get("https://alfabank.ru/ext-json/0.2/exchange/cash/")
	if err != nil {
		log.Fatalf("Can't get rates %v", err)
	}
	defer resp.Body.Close()
	var m AlfaExchange
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&m)
	if err != nil {
		log.Fatalf("Can't parse JSON %v", err)
	}
	var sellPrice float64
	var sellOrder string
	var buyPrice float64
	var buyOrder string
	for i:= range m.Response.Data.Usd {
		if m.Response.Data.Usd[i].Type == "sell" {
			sellPrice = m.Response.Data.Usd[i].Value
			sellOrder = m.Response.Data.Usd[i].Order
		}
		if m.Response.Data.Usd[i].Type == "buy" {
			buyPrice = m.Response.Data.Usd[i].Value
			buyOrder = m.Response.Data.Usd[i].Order
		}
	}
	fmt.Printf("%s%.2f %s%.2f", buyOrder, buyPrice, sellOrder, sellPrice)
}


