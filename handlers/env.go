package handlers

import (
	"fmt"
	"log"
	"os"
	"path"
)

type EnvRepairFlag int

const (
	NvimFolder EnvRepairFlag = iota
	EnvimLuaFile
	PluginsFolder
)

type CheckEnvironment struct {
	state  HandlerState
	rflags []EnvRepairFlag
	errors []error
	Path   string
}

func (ce *CheckEnvironment) GetType() HandlerType {
	return CheckEnvironmentType
}

func (ce *CheckEnvironment) GetState() HandlerState {
	return ce.state
}

func (ce *CheckEnvironment) GetErrors() []error {
	return ce.errors
}

func (ce *CheckEnvironment) Execute(state map[HandlerType]Handler) {
	if ce.state != HandlerNotStarted {
		log.Panic("Can't execute a handler that's already started.")
	}

	ce.state = HandlerError

	if ce.Path == "" {
		if p, err := os.Getwd(); err != nil {
			log.Panic("Error getting current working directory")
		} else {
			ce.Path = p
		}
	}

	if exist, creatable := statPath(path.Join(ce.Path, ".nvim"), true); !exist && creatable {
		ce.rflags = append(ce.rflags, NvimFolder)
	} else if !exist && !creatable {
		ce.errors = append(ce.errors, fmt.Errorf(".nvim is not creatable in %s", ce.Path))
	}

	if exist, creatable := statPath(path.Join(ce.Path, ".nvim", "envim.lua"), false); !exist && creatable {
		ce.rflags = append(ce.rflags, EnvimLuaFile)
	} else if !exist && !creatable {
		ce.errors = append(ce.errors, fmt.Errorf("envim.lua is not creatable in %s", path.Join(ce.Path, ".nvim")))
	}

	if exist, creatable := statPath(path.Join(ce.Path, ".nvim", "plugins"), true); !exist && creatable {
		ce.rflags = append(ce.rflags, PluginsFolder)
	} else if !exist && !creatable {
		ce.errors = append(ce.errors, fmt.Errorf("plugins directory is not creatable in %s", path.Join(ce.Path, ".nvim")))
	}

	if len(ce.errors) == 0 {
		ce.state = HandlerSuccess
	}
}

func (ce *CheckEnvironment) DependsOn() []HandlerType {
	return []HandlerType{}
}

func (ce *CheckEnvironment) ShouldProceed() bool {
	return ce.state != HandlerError
}

func (ce *CheckEnvironment) GetRepairFlags() []EnvRepairFlag {
	return ce.rflags
}
