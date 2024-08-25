package handlers

import (
	"log"
)

type HandlerType int

const (
	ChainType HandlerType = iota
	CheckEnvironmentType
	CreateEnvironmentType
)

var hTranslate = map[HandlerType]string{
	ChainType:             "Chain",
	CheckEnvironmentType:  "CheckEnvironment",
	CreateEnvironmentType: "CreateEnvironment",
}

func (h HandlerType) String() string {
	return hTranslate[h]
}

type HandlerState int

const (
	HandlerNotStartedState = iota
	HandlerErrorState
	HandlerSuccessState
)

type Handler interface {
	Execute(map[HandlerType]Handler)
	GetType() HandlerType
	ShouldProceed() bool
	GetState() HandlerState
	GetErrors() []error
	DependsOn() []HandlerType
}

func confirmExecution(h Handler) {
	if h.GetState() != HandlerNotStartedState {
		log.Panicf("%s: Cannot execute a handler that has already been executed", h.GetType())
	}
}

func GetHandler[T Handler](state map[HandlerType]Handler) T {
	var result T
  t := result.GetType()
	if h, ok := state[t]; ok {
		if result, ok = h.(T); !ok {
			log.Panicf("Something wierd happened: " + 
        "The %s in the execution state doesn't cast into its own type." +
        "Probably a duplicate in GetType.", t)
		}
	} else {
		log.Panicf("There is no %s in the current execution state.", t)
	}
	return result
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
	if c.state != HandlerNotStartedState {
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

	confirmExecution(c)

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
			c.state = HandlerErrorState
			return
		}

	}
	c.proceed = true
	c.state = HandlerSuccessState
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
