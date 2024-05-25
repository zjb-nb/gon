package middleware

import (
	"fmt"
	"gon/gonweb"
	"time"
)

type ComputeTimeSpendMiddleware struct {
	start int
	end   int
}

func (m *ComputeTimeSpendMiddleware) ComputerTimeSpend(next gonweb.GonHandlerFunc) gonweb.GonHandlerFunc {
	return func(c *gonweb.GonContext) {
		m.start = time.Now().Nanosecond()
		next(c)
		m.end = time.Now().Nanosecond()
		fmt.Printf("time spend %d ms\n", (m.end-m.start)/1000)
	}
}

func NewTimeMiddleware() *ComputeTimeSpendMiddleware {
	return &ComputeTimeSpendMiddleware{}
}

func SayHelloFilterBuilder(next gonweb.GonHandlerFunc) gonweb.GonHandlerFunc {
	return func(c *gonweb.GonContext) {
		fmt.Println("生命周期：开始")
		next(c)
		fmt.Println("生命周期：结束")
	}
}
