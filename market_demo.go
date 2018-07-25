package main

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"encoding/json"
	"time"
	"net/http"
	"io/ioutil"
)

type  CBDepthJson struct{
	Status		string	//	ok, error	
	Timestamp	int64	//		
	Symbol		string	//		
	OrderBook 	OderBookInfo
}
type OderBookInfo struct{
	Asks 		[]Book
	Bids 		[]Book
}
type Book struct{
	Price   	float64
    Quantity 	float64
}


type  CBTradesJson struct{
	Status		string	//	ok, error	
	Timestamp	int64	//		
	Symbol		string	//		
	Trades 		[]CBTrade
}

type CBTrade  struct{
	TradeId		string
	Price		string
	Quantity	string
	Take		string	// buy, sell	主动买入 / 主动卖出
	Time		int64
} 

type CBTickerJson struct{	
	Status		string	//	ok, error	
	Timestamp	int64	//		
	Ticker 		[]CBTicker
}
type CBTicker struct{
	Symbol		string//	true		
	Last		string//	true		
	Bid			string//	true		买一价
	Ask			string//	true		卖一价
	High24hr	string 	`json:"24hrHigh"`//	true		
	Low24hr		string 	`json:"24hrLow"`//	true		
	Vol24hr		string 	`json:"24hrVol"`//	true		
	Amt24hr		string 	`json:"24hrAmt"`//	true		
}

type SymbolInfo struct{
	Symbol  	string
	Depth 		string
	Size 		string
}

var (
	HttpUrl="http://api.coinbene.com/v1/market/"
)

func main() {


	testSymbol:=SymbolInfo{Symbol:"btcusdt",Depth:"20",Size:"10",}

	// 获取挂单
	go GoGetOrderBook(testSymbol)	

	// 获取成交记录
	go GoGetTrades(testSymbol)	

	// 获取最新价
	go GoGetTicker(testSymbol)
		

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

}


func GoGetOrderBook(symbol SymbolInfo) {
	for{
		time.Sleep(time.Millisecond * 1000)
		url := HttpUrl+"orderbook?symbol=" + symbol.Symbol+"&depth=" + symbol.Depth
		// fmt.Println(url)
		client := &http.Client{}
		reqest, err := http.NewRequest("GET", url, nil)

		if err != nil {
			fmt.Println("GoGetOrderBook http.NewRequest err:",err)
			continue
		}
		//处理返回结果
		response, err1 := client.Do(reqest)
		if err1 != nil {			
			fmt.Println("GoGetOrderBook client.Do err:",err1)
			continue
		}
		defer response.Body.Close()
		body, err2 := ioutil.ReadAll(response.Body)
		if err2 != nil {
			fmt.Println("GoGetOrderBook ioutil.ReadAll err", err2)
			continue
		}
		// fmt.Println(string(body))

		var getresult CBDepthJson
		jsonerr := json.Unmarshal(body, &getresult)
		if jsonerr != nil {
			fmt.Println("GoGetOrderBook result jsonerr:", jsonerr)
			continue
		} else {
			if getresult.Status == "ok" {
				fmt.Println("GoGetOrderBook ",getresult)
			}else{
				fmt.Println("GoGetOrderBook status err:", getresult.Status)
			}

		}

	}
	
}


func GoGetTrades(symbol SymbolInfo) {
	for{
		time.Sleep(time.Millisecond * 1000)
		url := HttpUrl+"trades?symbol=" + symbol.Symbol+"&size="+symbol.Size
		// fmt.Println(url)
		client := &http.Client{}
		reqest, err := http.NewRequest("GET", url, nil)

		if err != nil {
			fmt.Println("GoGetTrades http.NewRequest err:",err)
			continue
		}
		//处理返回结果
		response, err1 := client.Do(reqest)
		if err1 != nil {			
			fmt.Println("GoGetTrades client.Do err:",err1)
			continue
		}
		defer response.Body.Close()
		body, err2 := ioutil.ReadAll(response.Body)
		if err2 != nil {
			fmt.Println("GoGetTrades ioutil.ReadAll err", err2)
			continue
		}

		// fmt.Println(string(body))

		var getresult CBTradesJson
		jsonerr := json.Unmarshal(body, &getresult)
		if jsonerr != nil {
			fmt.Println("GoGetTrades result jsonerr:", jsonerr)			
			continue
		} else {
			if getresult.Status == "ok" {
				fmt.Println("GoGetTrades ",getresult)
			}else{
				fmt.Println("GoGetTrades status err:", getresult.Status)
			}
		}

	}
	
}

func GoGetTicker(symbol SymbolInfo) {
	for{
		time.Sleep(time.Millisecond * 1000)
		url := HttpUrl+"ticker?symbol=" + symbol.Symbol
		// fmt.Println(url)
		client := &http.Client{}
		reqest, err := http.NewRequest("GET", url, nil)

		if err != nil {
			fmt.Println("GoGetTicker http.NewRequest err:",err)
			continue
		}
		//处理返回结果
		response, err1 := client.Do(reqest)
		if err1 != nil {			
			fmt.Println("GoGetTicker client.Do err:",err1)
			continue
		}
		defer response.Body.Close()
		body, err2 := ioutil.ReadAll(response.Body)
		if err2 != nil {
			fmt.Println("GoGetTicker ioutil.ReadAll err", err2)
			continue
		}


		var getresult CBTickerJson
		jsonerr := json.Unmarshal(body, &getresult)
		if jsonerr != nil {
			fmt.Println("GoGetTicker result jsonerr:", jsonerr)
			fmt.Println(string(body))
			continue
		} else {
			if getresult.Status == "ok" {
				fmt.Println("GoGetTicker ",getresult)
			}else{
				fmt.Println("GoGetTicker status err:", getresult.Status)
			}
		}

	}
	
}

