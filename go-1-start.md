1.web的体系：
	基于http协议，请求响应模式。
	对应一个web应用，主要有这几点：
		(1)监听端口
		(2)路由规则映射
			
2.go语言生成web服务

	go内置了net/http标准库跑web,test.go，cmd执行go run test.go,浏览器打开：http://localhost，http://localhost/json,发出请求，然后可以看到响应结果。
	
	访问过程：
		服务端监听80端口，服务启动时候，路由规则由map管理：映射"/"和"/json"
		客户端请求http://localhost/json，host(使用默认端口80)为：localhost，path："/json"
	package main

	//导包
	import (
		"net/http"  //http协议
		"encoding/json"  //json解析
		"fmt"  //读写
	)

	func main()  {
		//返回hello world
		http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprint(writer,"hello world")
		})
		//返回json
		http.HandleFunc("/json", func(writer http.ResponseWriter, request *http.Request) {
			word,_:=json.Marshal(map[string]string{"word":"hello world"})
			fmt.Fprint(writer,string(word))
		})
		http.ListenAndServe(":80",nil)
	}



