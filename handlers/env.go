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

var rfTranslate = map[EnvRepairFlag]string{
	NvimFolder:    "NvimFolder",
	EnvimLuaFile:  "EnvimLuaFile",
	PluginsFolder: "PluginsFolder",
}

func (rf EnvRepairFlag) String() string {
	return rfTranslate[rf]
}

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
	confirmExecution(ce)

	ce.state = HandlerErrorState

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
		ce.state = HandlerSuccessState
	}
}

func (ce *CheckEnvironment) DependsOn() []HandlerType {
	return []HandlerType{}
}

func (ce *CheckEnvironment) ShouldProceed() bool {
	return ce.state == HandlerSuccessState
}

func (ce *CheckEnvironment) GetRepairFlags() []EnvRepairFlag {
	return ce.rflags
}

// CreateEnvironment creates an environment in a given path
type CreateEnvironment struct {
	state  HandlerState
	errors []error
	Path   string
}

func (ce *CreateEnvironment) GetType() HandlerType {
	return CreateEnvironmentType
}

func (ce *CreateEnvironment) GetState() HandlerState {
	return ce.state
}

func (ce *CreateEnvironment) ShouldProceed() bool {
	return ce.state == HandlerSuccessState
}

func (ce *CreateEnvironment) GetErrors() []error {
	return ce.errors
}

func (ce *CreateEnvironment) DependsOn() []HandlerType {
	return []HandlerType{CheckEnvironmentType}
}

func (ce *CreateEnvironment) Execute(state map[HandlerType]Handler) {
	confirmExecution(ce)

	ce.state = HandlerErrorState

	check := GetHandler[*CheckEnvironment](state)

	if ce.Path == "" {
		if p, err := os.Getwd(); err != nil {
			log.Panic("Error getting current working directory")
		} else {
			ce.Path = p
		}
	}

	for _, rflag := range check.GetRepairFlags() {
		switch rflag {
		case NvimFolder:
			rpath := ".nvim"
			path := path.Join(ce.Path, rpath)
			if created, err := createFolder(path); !created {
				ce.errors = append(ce.errors, err)
				log.Printf("Folder %s already exists in %s. Skipping...\n", rpath, ce.Path)
			} else if err != nil {
				log.Printf("Folder %s already exists in %s. Skipping...\n", rpath, ce.Path)
			}
		case EnvimLuaFile:
		case PluginsFolder:
			rpath := path.Join(".nvim", "plugins")
			path := path.Join(ce.Path, rpath)
			if created, err := createFolder(path); !created {
				ce.errors = append(ce.errors, err)
				log.Printf("Folder %s already exists in %s. Skipping...\n", rpath, ce.Path)
			} else if err != nil {
				log.Printf("Folder %s already exists in %s. Skipping...\n", rpath, ce.Path)
			}
		default:
			log.Panicf("%s: There is no implementation for %s EnvRepairFlag", ce.GetType(), rflag)
		}
	}

	ce.state = HandlerSuccessState
}
