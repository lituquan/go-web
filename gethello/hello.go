package main

//导包
import (
	"net/http"  //http协议
	"fmt"  //读写
	"io/ioutil"
)

func httpGet()string {
    resp, err := http.Get("http://106.12.27.170/json")
    if err != nil {
        // handle error
		return "Get请求失败"
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        // handle error
		return "Get请求读取失败"
    }
    fmt.Println(string(body))
	return string(body)
}
func main()  {	
	//返回json
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		word:=httpGet()
		fmt.Fprint(writer,word)
	})
	http.ListenAndServe(":8080",nil)
}