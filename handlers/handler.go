package handlers

import (
	"log"
)

type HandlerType int

const (
	ChainType HandlerType = iota
	CheckEnvironmentType
)

var hTranslate = map[HandlerType]string{
  ChainType: "Chain",
  CheckEnvironmentType: "CheckEnvironment",
}

func (h HandlerType) String() string {
  return hTranslate[h]
}

type HandlerState int

const (
	HandlerNotStarted = iota
	HandlerError
	HandlerSuccess
)

type Handler interface {
	Execute(map[HandlerType]Handler)
	GetType() HandlerType
	ShouldProceed() bool
	GetState() HandlerState
	GetErrors() []error
  DependsOn() []HandlerType
}

// Chain Handler is a chain of multiple handlers that are executed in order.
// Chain Handler is a self sufficient unit of code, meaning that it can't contain
// any undetermined dependencies. 

type ChainHandler struct {
	chain   []Handler
	errors  []error
	exeMap  map[HandlerType]Handler
	state   HandlerState
	proceed bool
}

func (c *ChainHandler) GetType() HandlerType {
	return ChainType
}

func (c *ChainHandler) ShouldProceed() bool {
	return c.proceed
}

func (c *ChainHandler) GetState() HandlerState {
	return c.state
}

func (c *ChainHandler) GetErrors() []error {
	return c.errors
}

func (c *ChainHandler) AddHandler(h Handler) {
	if c.state != HandlerNotStarted {
		log.Panic("Cannot add a handler to a chain that has already been executed.")
	}

  for _, d := range h.DependsOn() {
    if _, ok := c.exeMap[d]; !ok {
      log.Panicf("Handler %s depends on %s, but %s is not in the chain.", h.GetType(), d, d)
    }
  }
	c.chain = append(c.chain, h)
}

func (c *ChainHandler) Execute(state map[HandlerType]Handler) {

	c.exeMap = make(map[HandlerType]Handler)

	if state != nil {
		c.exeMap = state
	}

	for _, h := range c.chain {
		h.Execute(c.exeMap)
    c.exeMap[h.GetType()] = h
		c.errors = append(c.errors, h.GetErrors()...)

		if !h.ShouldProceed() {
			c.proceed = false
			c.state = HandlerError
			return
		}

	}
	c.proceed = true
	c.state = HandlerSuccess
}

func (c *ChainHandler) DependsOn() []HandlerType {
  return []HandlerType{}
}

func (c *ChainHandler) GetHandler(t HandlerType) Handler {
  if h, ok := c.exeMap[t]; ok {
    return h
  } else {
    return nil
  }
}
