package servant

import (
	"context"
	"sync"
	"time"
)

type Context struct {
	context.Context
	mu           sync.RWMutex
	param        map[string]any
	index        int
	handlerChain HandlersChain
}

func (c *Context) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists = c.param[key]

	return
}

func (c *Context) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.param == nil {
		c.param = make(map[string]any)
	}

	c.param[key] = value
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlerChain) {
		c.handlerChain[c.index](c)
		c.index++
	}
}

/************************************/
/***** GOLANG.ORG/X/NET/CONTEXT *****/
/************************************/

// Deadline always returns that there is no deadline (ok==false),
// maybe you want to use Request.Context().Deadline() instead.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done always returns nil (chan which will wait forever),
// if you want to abort your work when the connection was closed
// you should use Request.Context().Done() instead.
func (c *Context) Done() <-chan struct{} {
	return nil
}

// Err always returns nil, maybe you want to use Request.Context().Err() instead.
func (c *Context) Err() error {
	return nil
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *Context) Value(key any) any {
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return c.Context.Value(key)
}
