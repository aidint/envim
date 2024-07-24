package handler

import (
	"log"
)

type HandlerType int

const (
	ChainHandlerType = iota
)

type Handler interface {
	Execute(map[HandlerType]*Handler)
	GetType() HandlerType
  ShouldContinue() bool
  GetErrors() []error
}

type ChainHandler struct {
  chain []*Handler
  shouldContinue bool
  errors []error
  exeMap map[HandlerType]*Handler
  Done bool
}

func (c *ChainHandler) GetType() HandlerType {
  return ChainHandlerType
}

func (c *ChainHandler) ShouldContinue() bool {
  return c.shouldContinue
}

func (c *ChainHandler) GetErrors() []error {
  return c.errors
}

func (c *ChainHandler) AddHandler(h *Handler) {
  if c.Done {
    log.Panic("Cannot add a handler to a chain that has already been executed.")
  }
  
  c.chain = append(c.chain, h)
  c.shouldContinue = true
}

func (c *ChainHandler) Execute(state map[HandlerType]*Handler) {
  if c.Done {
    log.Panic("Cannot execute a chain that has already been executed.")
  }
  c.exeMap = state
}
