package kernel

import (
	"github.com/firmeve/firmeve/binding"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/render"
	"time"
)

const (
	abortIndex = -1
)

type (
	context struct {
		firmeve  contract.Application
		protocol contract.Protocol
		handlers []contract.ContextHandler
		entries  map[string]*contract.ContextEntity
		index    int
	}
)

func NewContext(firmeve contract.Application, protocol contract.Protocol, handlers ...contract.ContextHandler) contract.Context {
	return &context{
		firmeve:  firmeve,
		protocol: protocol,
		handlers: handlers,
		entries:  make(map[string]*contract.ContextEntity, 0),
		index:    0,
	}
}

func (c *context) Firmeve() contract.Application {
	return c.firmeve
}

func (c *context) Protocol() contract.Protocol {
	return c.protocol
}

func (c *context) Error(status int, err error) {
	var newErr contract.ErrorRender
	if v, ok := err.(contract.ErrorRender); ok {
		newErr = v
	} else {
		newErr = ErrorWarp(err)
	}

	if err2 := newErr.Render(status, c); err2 != nil {
		panic(err2)
	}
}

func (c *context) Abort() {
	c.index = abortIndex
}

func (c *context) Next() {
	if c.index == abortIndex {
		return
	}

	if c.index < len(c.handlers) {
		c.index++
		c.handlers[c.index-1](c)
	}
}

func (c *context) Handlers() []contract.ContextHandler {
	return c.handlers
}

func (c *context) AddEntity(key string, value interface{}) {
	c.entries[key] = &contract.ContextEntity{
		Key:   key,
		Value: value,
	}
}

func (c *context) Entity(key string) *contract.ContextEntity {
	if v, ok := c.entries[key]; ok {
		return v
	}

	return nil
}

func (c *context) Get(key string) interface{} {
	values := c.protocol.Values()
	if value, ok := values[key]; ok {
		return value
	}

	return nil
}

func (c *context) Bind(v interface{}) error {
	return binding.Bind(c.protocol, v)
}

func (c *context) BindWith(b contract.Binding, v interface{}) error {
	return b.Protocol(c.protocol, v)
}

func (c *context) RenderWith(status int, r contract.Render, v interface{}) error {
	return r.Render(c.protocol, status, v)
}

func (c *context) Render(status int, v interface{}) error {
	return render.Render(c.protocol, status, v)
}

func (c *context) Clone() contract.Context {
	//@todo 暂时先返回自己，Context全部完善后再修改clone
	return c
}

// --------------------------- context.Context -> Base context ------------------------

func (c *context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *context) Done() <-chan struct{} {
	return nil
}

func (c *context) Err() error {
	return nil
}

func (c *context) Value(key interface{}) interface{} {
	if v, ok := key.(string); ok {
		return c.Get(v)
	}

	return nil
}
