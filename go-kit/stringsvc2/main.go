package main

import (
	"net/http"
	"os"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	//日志中间件
	logger := log.NewLogfmtLogger(os.Stderr)

	//监控指标
	fieldKeys := []string{"method", "error"}
	/*
	type Counter interface {
		With(labelValues ...string) Counter
		Add(delta float64)
	}
	*/
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",//全部请求数目
		Help:      "Number of requests received.",
	}, fieldKeys)
	/*
	type Histogram interface { //直方图
		With(labelValues ...string) Histogram
		Observe(value float64)
	}
	*/
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",//请求数/微秒
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",//每个方法请求数目
	}, []string{}) // no fields here

	var svc StringService
	svc = stringService{}
	svc = loggingMiddleware{logger, svc}//日志
	svc = instrumentingMiddleware{requestCount, requestLatency, countResult, svc}
	//顺序:dec--svc:endpoint--logger--requestCount--requestLatency--countResult--env
	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)
	//注册路由和控制器
	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", ":8080")
	//开启服务
	logger.Log("err", http.ListenAndServe(":8080", nil))
}
