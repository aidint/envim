package handlers

import (
	"log"
)

type HandlerType int

const (
	ChainType = iota
  StatPathType
	CreateFolderType
  CheckEnvironmentType
)

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
}

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

	c.chain = append(c.chain, h)
}

func (c *ChainHandler) Execute(state map[HandlerType]Handler) {
	if c.state != HandlerNotStarted {
		log.Panic("Cannot execute a chain that has already been executed.")
	}

	c.exeMap = make(map[HandlerType]Handler)

	if state != nil {
		c.exeMap = state
	}

	for _, h := range c.chain {
		h.Execute(c.exeMap)
		c.errors = append(c.errors, h.GetErrors()...)

		if !h.ShouldProceed() {
			c.proceed = false
			c.state = HandlerError
      return
		}

		c.exeMap[h.GetType()] = h
	}
	c.proceed = true
	c.state = HandlerSuccess
}

