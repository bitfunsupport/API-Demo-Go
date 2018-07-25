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
	"strings"
	"strconv"
	"crypto/md5"
	"encoding/hex"
	"sort"
)

type AddOrderInfo struct{
	Symbol 		string
	Type  		string
	Price 		string
	Quantity 		string
}


var (

	AccId = "accid"
	AccSecret = "accsecret"

	BalanceUrl = "http://api.coinbene.com/v1/trade/balance"
	AddOrderUrl = "http://api.coinbene.com/v1/trade/order/place"
	CancelOrderUrl = "http://api.coinbene.com/v1/trade/order/cancel"
	CheckOrderUrl = "http://api.coinbene.com/v1/trade/order/info"
	OrderListUrl = "http://api.coinbene.com/v1/trade/order/open-orders"
)


type BalanceResult struct{
	Status			string	//true	ok, error	
	Account 		string
	Timestamp		int64 //long	true		
	Balance 	 	[]BlanceInfo
}

type BlanceInfo struct{
 	Asset		string	//		资产名称/缩写
	Available	string	//		可用余额
	Reserved	string	//		冻结余额
	Total 		string  //  	总数
}


type AddOrderResult struct {
	Status 		string
	Timestamp    int64
	OrderId   	string
}

type CancelOroderResult struct {
	Status 		string
	Timestamp    int64
	OrderId   	string
}

type CheckOrderResult struct{
	Status			string	//true	ok, error	
	Symbol			string	//true		
	Timestamp		int64 //long	true		
	Order 	 		OrderInfo
}

type  OrderInfo struct{	
	OrderId 		string
	OrderStatus 	string
	Symbol 			string
	Type  			string
	Price 			string
	OrderQuantity 	string
	FilledQuantity	string	
	FilledAmount	string
	AveragePrice 	string
	Fees			string
	LastModified 	string
	CreatedTime		string
}


type OrderListResult struct{
	Status			string	//true	ok, error
	Timestamp		int64 //long	true		
	Orders 	 	[]OrderInfo
}



func main() {

	go GetAccountBalance()

	var testorder AddOrderInfo
	testorder.Symbol="ethusdt"
	testorder.Type="sell-limit"
	testorder.Price="520"
	testorder.Quantity="1.12"
	//go 	AddOrder(testorder)
	

	//go CancelOrder("orderid")

	go CheckOrder("testorderid")

	go GetOrderList("btcusdt")
	

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

}


func GetAccountBalance() {
	timestampstr := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)	
	parameters := make(map[string]string)
	parameters["account"]="exchange"
	parameters["apiid"]=AccId
	parameters["secret"]=AccSecret
	parameters["timestamp"]=timestampstr
	postbodystr:=GetPostBodyString(parameters)

	client := &http.Client{}
	reqest, err := http.NewRequest("POST", BalanceUrl, strings.NewReader(postbodystr))
	if err != nil {
		fmt.Println("GetAccountBalance http.NewRequest err",err)
		return
	}
	reqest.Header.Set("Content-Type", "application/json; charset=utf-8")
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko")
	reqest.Header.Set("Connection", "keep-alive")
	//处理返回结果
	response, errdo := client.Do(reqest)
	if errdo != nil {
		fmt.Println("GetAccountBalance client.Do err",errdo)
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("GetAccountBalance ioutil.ReadAll err", err)
		return
	}
	// fmt.Println(timestampstr,string(body))
	

	var balanceresult BalanceResult
	jsonerr := json.Unmarshal(body, &balanceresult)
	if jsonerr != nil {
		fmt.Println("GetAccountBalance result jsonerr:", jsonerr)		
		return
	} else {
		if balanceresult.Status == "ok" {
			fmt.Println("GetAccountBalance result:",balanceresult)
		}else{
			fmt.Println("GetAccountBalance result err:",balanceresult.Status,string(body))
		}

	}

}

func AddOrder(orderinfo AddOrderInfo) {
	timestampstr := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	parameters := make(map[string]string)
	parameters["apiid"]=AccId
	parameters["secret"]=AccSecret
	parameters["timestamp"]=timestampstr
	parameters["price"]=orderinfo.Price
	parameters["quantity"]=orderinfo.Quantity
	parameters["symbol"]=orderinfo.Symbol
	parameters["type"]=orderinfo.Type
	postbodystr:=GetPostBodyString(parameters)

	client := &http.Client{}
	reqest, err := http.NewRequest("POST", AddOrderUrl, strings.NewReader(postbodystr))
	if err != nil {
		fmt.Println("AddOrder http.NewRequest err",err)
		return
	}
	reqest.Header.Set("Content-Type", "application/json; charset=utf-8")
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko")
	reqest.Header.Set("Connection", "keep-alive")

	//处理返回结果
	response, posterr := client.Do(reqest)
	if posterr != nil {
		fmt.Println("AddOrder post err ", posterr)		
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err", err)		
		return
	}
	// fmt.Println(timestampstr,string(body))

	var addresult AddOrderResult
	jsonerr := json.Unmarshal(body, &addresult)
	if jsonerr != nil {
		fmt.Println("AddOrder result jsonerr:", jsonerr)		
		return
	} else {
		if addresult.Status == "ok" {
			fmt.Println("AddOrder result:",addresult)
		}else{
			fmt.Println("AddOrder result err:",addresult.Status,string(body))
		}

	}

}


func CancelOrder(orderid string) {
	timestampstr := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	parameters := make(map[string]string)
	parameters["apiid"]=AccId
	parameters["secret"]=AccSecret
	parameters["timestamp"]=timestampstr
	parameters["orderid"]=orderid
	postbodystr:=GetPostBodyString(parameters)
	
	client := &http.Client{}
	reqest, err := http.NewRequest("POST", CancelOrderUrl, strings.NewReader(postbodystr))
	if err != nil {
		fmt.Println("CancelOrder http.NewRequest err",err)
		return
	}
	reqest.Header.Set("Content-Type", "application/json; charset=utf-8")
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko")
	reqest.Header.Set("Connection", "keep-alive")

	//处理返回结果
	response, posterr := client.Do(reqest)
	if posterr != nil {
		fmt.Println("CancelOrder post err ", posterr)		
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err", err)		
		return
	}
	// fmt.Println(timestampstr,string(body))

	var cancelresult CancelOroderResult
	jsonerr := json.Unmarshal(body, &cancelresult)
	if jsonerr != nil {
		fmt.Println("CancelOrder result jsonerr:", jsonerr)		
		return
	} else {
		if cancelresult.Status == "ok" {
			fmt.Println("CancelOrder result:",cancelresult)
		}else{
			fmt.Println("CancelOrder result err:",cancelresult.Status,string(body))
		}

	}

}

func CheckOrder(orderid string) {
	timestampstr := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	parameters := make(map[string]string)
	parameters["apiid"]=AccId
	parameters["secret"]=AccSecret
	parameters["timestamp"]=timestampstr
	parameters["orderid"]=orderid
	postbodystr:=GetPostBodyString(parameters)

	client := &http.Client{}
	reqest, err := http.NewRequest("POST", CheckOrderUrl, strings.NewReader(postbodystr))
	if err != nil {
		fmt.Println("CheckOrder http.NewRequest err",err)
		return
	}
	reqest.Header.Set("Content-Type", "application/json; charset=utf-8")
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko")
	reqest.Header.Set("Connection", "keep-alive")

	//处理返回结果
	response, posterr := client.Do(reqest)
	if posterr != nil {
		fmt.Println("CheckOrder post err ", posterr)		
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err", err)		
		return
	}
	// fmt.Println(timestampstr,string(body))

	var chekcresult CheckOrderResult
	jsonerr := json.Unmarshal(body, &chekcresult)
	if jsonerr != nil {
		fmt.Println("CheckOrder result jsonerr:", jsonerr)		
		return
	} else {
		if chekcresult.Status == "ok" {
			fmt.Println("CheckOrder result:",chekcresult)
		}else{
			fmt.Println("CheckOrder result err:",chekcresult.Status,string(body))
		}

	}

}


func GetOrderList(symbol string) {
	timestampstr := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	parameters := make(map[string]string)
	parameters["apiid"]=AccId
	parameters["secret"]=AccSecret
	parameters["timestamp"]=timestampstr
	parameters["symbol"]=symbol
	postbodystr:=GetPostBodyString(parameters)


	client := &http.Client{}
	reqest, err := http.NewRequest("POST", OrderListUrl, strings.NewReader(postbodystr))
	if err != nil {
		fmt.Println("GetOrderList http.NewRequest err",err)
		return
	}
	reqest.Header.Set("Content-Type", "application/json; charset=utf-8")
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko")
	reqest.Header.Set("Connection", "keep-alive")
	//处理返回结果
	response, errdo := client.Do(reqest)
	if errdo != nil {
		fmt.Println("GetOrderList client.Do err",errdo)
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("GetOrderList ioutil.ReadAll err", err)
		return
	}
	// fmt.Println(timestampstr,string(body))
	

	var orderlistresult OrderListResult
	jsonerr := json.Unmarshal(body, &orderlistresult)
	if jsonerr != nil {
		fmt.Println("GetOrderList result jsonerr:", jsonerr)		
		return
	} else {
		if orderlistresult.Status == "ok" {
			fmt.Println("GetOrderList result:",orderlistresult)
		}else{
			fmt.Println("GetOrderList result err:",orderlistresult.Status,string(body))
		}

	}

}






func GetPostBodyString(parameters map[string]string)(postbody string){
	var keys []string
    for k := range parameters {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    unsignstr := ""
    for _, k := range keys {
    	unsignstr=unsignstr+k+"="+parameters[k]+"&"
    }
    unsignstr=strings.ToUpper(unsignstr[0:len(unsignstr)-1])
    signstr:=GetMD5Str(unsignstr)
    delete(parameters,"secret")
    parameters["sign"] = signstr

    for key,value := range parameters{
    	postbody = postbody + `"` + key + `":"` + value + `",`
    }
    postbody = `{`+ postbody[0:len(postbody)-1] + `}`
    return
}


func GetMD5Str(str string) string{
	ret:=""
	md5Ctx := md5.New()
    md5Ctx.Write([]byte(str))
    cipherStr := md5Ctx.Sum(nil)
    ret = hex.EncodeToString(cipherStr)
	return ret
}