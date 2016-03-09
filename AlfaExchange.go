package main

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
)

//type AlfaFilter struct {
//	text string
//	segment string
//}
//type AlfaRequest struct {
//	server string
//	service string
//	version string
//	order string
//	limit int
//	offset int
//	filter AlfaFilter
//}
//
//type AlfaResponse struct {
//	status string
//}
//
//type AlfaExchange struct {
//	request AlfaRequest
//	response AlfaResponse
//}

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
				     Chf []struct {
					     Date  string  `json:"date"`
					     Order string  `json:"order"`
					     Type  string  `json:"type"`
					     Value float64 `json:"value"`
				     } `json:"chf"`
				     Eur []struct {
					     Date  string  `json:"date"`
					     Order string  `json:"order"`
					     Type  string  `json:"type"`
					     Value float64 `json:"value"`
				     } `json:"eur"`
				     Gbp []struct {
					     Date  string  `json:"date"`
					     Order string  `json:"order"`
					     Type  string  `json:"type"`
					     Value float64 `json:"value"`
				     } `json:"gbp"`
				     Usd []struct {
					     Date  string  `json:"date"`
					     Order string  `json:"order"`
					     Type  string  `json:"type"`
					     Value float64 `json:"value"`
				     } `json:"usd"`
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
	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	var m AlfaExchange
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&m)
	//err = json.Unmarshal([]byte(body), &m)
	if err != nil {
		log.Fatalf("Can't parse JSON %v", err)
	}
	fmt.Printf("%s%.2f %s%.2f", m.Response.Data.Usd[0].Order,m.Response.Data.Usd[0].Value,
		m.Response.Data.Usd[1].Order,m.Response.Data.Usd[1].Value)
}


