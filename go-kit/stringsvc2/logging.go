package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger //挂上日志监控
	next   StringService //业务服务对象:组合，装饰模式增强对象
}

func (mw loggingMiddleware) Uppercase(s string) (output string, err error) {
	nanosecond("upper log")
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())//执行完方法就输出日志

	output, err = mw.next.Uppercase(s)
	return
}

func (mw loggingMiddleware) Count(s string) (n int) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())//执行完方法就输出日志

	n = mw.next.Count(s)
	return
}
