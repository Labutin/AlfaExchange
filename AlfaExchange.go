package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"
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

func getStoreFilePath() (dir string, fileName string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir + "/.AlfaExchange/", "minimum.txt"
}

func createStoreFileIfNotExist(storePath string, fileName string) {
	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		os.Mkdir(storePath, 0700)
	}
	if _, err := os.Stat(storePath + fileName); os.IsNotExist(err) {
		initValue := []byte("100000")
		err := ioutil.WriteFile(storePath+fileName, initValue, 0600)
		if err != nil {
			log.Fatalf("Can't create store file %v", err)
		}
	}
}

func getSeenMinimum() float64 {
	storePath, storeFileName := getStoreFilePath()

	createStoreFileIfNotExist(storePath, storeFileName)

	dat, err := ioutil.ReadFile(storePath + storeFileName)
	if err != nil {
		log.Fatalf("Can't read store file %v", err)
	}
	res, err := strconv.ParseFloat(string(dat), 64)
	if err != nil {
		log.Fatalf("Can't convert data from store file %v", err)
	}

	return res
}

func updateSeenMinimum(newMinimum float64) {
	initValue := []byte(strconv.FormatFloat(newMinimum, 'f', 2, 64))
	storePath, storeFileName := getStoreFilePath()
	err := ioutil.WriteFile(storePath+storeFileName, initValue, 0644)
	if err != nil {
		log.Fatalf("Can't write to store file %v", err)
	}
}

func main() {
	http.DefaultClient.Transport = &http.Transport{
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
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
	for i := range m.Response.Data.Usd {
		if m.Response.Data.Usd[i].Type == "sell" {
			sellPrice = m.Response.Data.Usd[i].Value
			sellOrder = m.Response.Data.Usd[i].Order
		}
		if m.Response.Data.Usd[i].Type == "buy" {
			buyPrice = m.Response.Data.Usd[i].Value
			buyOrder = m.Response.Data.Usd[i].Order
		}
	}

	if buyOrder == "0" {
		buyOrder = ""
	}
	if sellOrder == "0" {
		sellOrder = ""
	}
	outputFormat := "%s%.2f %s%.2f"
	seenMinimum := getSeenMinimum()
	if sellPrice <= seenMinimum {
		outputFormat = "%s%.2f %s%.2f | color=green"
		updateSeenMinimum(sellPrice)
	}
	fmt.Printf(outputFormat, buyOrder, buyPrice, sellOrder, sellPrice)
}
