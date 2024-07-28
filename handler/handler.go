package handler

import (
	"fmt"
	"log"
	"os"
)

type HandlerType int

const (
	ChainHandlerType = iota
	CreateFolderHandlerType
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
	return ChainHandlerType
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

// CreateFolderHandler

type CreateFolder struct {
	state      HandlerState
	errors     []error
	FolderName string
}

func (cf *CreateFolder) GetType() HandlerType {
	return CreateFolderHandlerType
}

func (cf *CreateFolder) GetState() HandlerState {
	return cf.state
}

func (cf *CreateFolder) GetErrors() []error {
	return cf.errors
}

func (cf *CreateFolder) Execute(state map[HandlerType]Handler) {
	if cf.state != HandlerNotStarted {
		log.Panic("Cannot execute a handler that has already been executed.")
	}

	if info, err := os.Stat(cf.FolderName); os.IsNotExist(err) {
		if err := os.MkdirAll(cf.FolderName, 0755); err != nil {
			cf.errors = append(cf.errors, fmt.Errorf("Create Folder %s error: %s", cf.FolderName, err.Error()))
			cf.state = HandlerError
			return
		}
		cf.state = HandlerSuccess
		return
	} else if info.IsDir() {
		cf.errors = append(cf.errors, fmt.Errorf("Create Folder %s error: %s", cf.FolderName, "Folder already exists"))
		cf.state = HandlerSuccess
		return
	}
	cf.state = HandlerError
	cf.errors = append(cf.errors, fmt.Errorf("Create Folder %s error: A file exists with the same name", cf.FolderName))
}

func (cf *CreateFolder) ShouldProceed() bool {
	return cf.state == HandlerSuccess
}
