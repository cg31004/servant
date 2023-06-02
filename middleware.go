package servant

import (
	"container/list"
)

var handler Handler

type HandleFunc func(ctx *CronContext)
type HandlersChain []HandleFunc

type Handler struct {
}

// AdapterHandle 可在主邏輯前後包 middleware
// mainHandler, middlewareA, middlewareB, middlewareC 處理順序
// middlewareA -> middlewareB -> middlewareC -> mainHandler -> middlewareC -> middlewareB -> middlewareA
func AdapterHandle(funcJob FuncJob, middleware ...HandleFunc) FuncJob {
	return handler.AdapterHandle(funcJob, middleware...)
}
func (h Handler) AdapterHandle(funcJob FuncJob, middleware ...HandleFunc) FuncJob {
	handles := append(middleware, h.convert(funcJob)...)
	return func(ctx *CronContext) {
		handlesChain := h.chain(handles...)
		c := &adapterContext{
			ctx:        ctx,
			middleware: handlesChain,
		}
		c.Next()
	}
}

// convert 把 mvc.RouteHandler 轉換成 HandleFunc
func (h Handler) convert(handlers ...FuncJob) []HandleFunc {
	length := len(handlers)
	result := make([]HandleFunc, length)
	for i, handler := range handlers {
		result[i] = func(ctx *CronContext) {
			handler(ctx)
		}
	}

	return result
}

// chain 把 HandleFunc 串起來形成 middleware 的串接處理流程
func (h Handler) chain(handlers ...HandleFunc) *list.List {
	l := list.New()

	for _, handler := range handlers {
		l.PushBack(handler)
	}

	return l
}

type AdapterContext interface {
	Next()
}

type adapterContext struct {
	ctx        *CronContext
	middleware *list.List
}

// Next 呼叫下一個 middleware 。
// 表示目前的 middleware 處理告一段落，先讓下一個 middleware 接著執行
// 如果沒有呼叫 Next() 表示為 middleware 串的最後一個方法，接著執行前面的 middleware 在呼叫 Next() 之後的邏輯
func (c *adapterContext) Next() {
	next := c.middleware.Remove(c.middleware.Front())
	if f, ok := next.(HandleFunc); ok {
		f(c.ctx)
	}
}
