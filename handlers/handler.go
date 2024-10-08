package handlers

import (
	"fmt"
	"log"
)

type HandlerType int

const (
	ChainType HandlerType = iota
	CheckEnvironmentType
	RepairEnvironmentType
)

var hTranslate = map[HandlerType]string{
	ChainType:             "Chain",
	CheckEnvironmentType:  "CheckEnvironment",
	RepairEnvironmentType: "RepairEnvironment",
}

func (h HandlerType) String() string {
	if h > 0 {
		return hTranslate[h]
	}
	return fmt.Sprintf("Mock %d", -1*h)
}

type HandlerState int

const (
	HandlerNotStartedState HandlerState = iota
	HandlerErrorState
	HandlerSuccessState
)

func (s HandlerState) String() string {
	switch s {
	case HandlerNotStartedState:
		return "Not Started"
	case HandlerErrorState:
		return "Error"
	case HandlerSuccessState:
		return "Success"
	default:
		return "Unknown"
	}
}

type Handler interface {
	Execute(map[HandlerType]Handler)
	GetType() HandlerType
	ShouldProceed() bool
	GetState() HandlerState
	GetErrors() []error
	DependsOn() []HandlerType
}

func prepareExecution(h Handler) {
	if h.GetState() != HandlerNotStartedState {
		log.Panicf("%s: Cannot execute a handler that has already been executed", h.GetType())
	}
}

func GetHandler[T Handler](state map[HandlerType]Handler) T {
	var result T
	t := result.GetType()
	if h, ok := state[t]; ok {
		if result, ok = h.(T); !ok {
			log.Panicf("Something wierd happened: "+
				"The %s in the execution state doesn't cast into its own type."+
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
	chain      []Handler
	errors     []error
	ExeMap     map[HandlerType]Handler
	state      HandlerState
	proceed    bool
	registered map[HandlerType]bool
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

	if c.registered == nil {
		c.registered = make(map[HandlerType]bool)
	}

	for _, d := range h.DependsOn() {
		if _, ok := c.registered[d]; !ok {
			if _, ok := c.ExeMap[d]; !ok {
				log.Panicf("Handler %s depends on %s, but %s is not in the chain.", h.GetType(), d, d)
			}
		}
	}

	c.chain = append(c.chain, h)
	c.registered[h.GetType()] = true
}

func (c *ChainHandler) Execute(state map[HandlerType]Handler) {

	prepareExecution(c)

	c.ExeMap = make(map[HandlerType]Handler)

	for _, h := range c.chain {
		h.Execute(c.ExeMap)
		c.ExeMap[h.GetType()] = h
		c.errors = append(c.errors, h.GetErrors()...)

		if !h.ShouldProceed() {
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
