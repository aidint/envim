package handlers

import (
	"envim/assets"
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

/*
CheckEnvironment checks if the environment is setup correctly.
We can set the path to check the environment in a different directory,
but by default it checks the current working directory.
It collects a list of repair flags defined earlier as EnvRepairFlag type.
It also collects a list of errors if any are encountered.
Based on the errors it'll change the handler state,
if handler state is not `HandlerSuccessState` it will send a stop signal
to the chain handler running it.
*/

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
	prepareExecution(ce)

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

/*
RepairEnvironment handler repairs the envim environment in the given path.
It depends on CheckEnvironment and therefore takes the path from it.
*/
type RepairEnvironment struct {
	state  HandlerState
	errors []error
}

func (re *RepairEnvironment) GetType() HandlerType {
	return RepairEnvironmentType
}

func (re *RepairEnvironment) GetState() HandlerState {
	return re.state
}

func (re *RepairEnvironment) ShouldProceed() bool {
	return re.state == HandlerSuccessState
}

func (re *RepairEnvironment) GetErrors() []error {
	return re.errors
}

func (re *RepairEnvironment) DependsOn() []HandlerType {
	return []HandlerType{CheckEnvironmentType}
}

func createConfigFile(filePath string) error {
	if err := os.WriteFile(filePath, []byte(assets.SampleConfig), 0644); err != nil {
		return err
	}
	return nil
}


func (re *RepairEnvironment) Execute(state map[HandlerType]Handler) {
	prepareExecution(re)

	re.state = HandlerErrorState

	ce := GetHandler[*CheckEnvironment](state)

	for _, rflag := range ce.GetRepairFlags() {
    
		switch rflag {

		case NvimFolder:
			rpath := ".nvim"
			p := path.Join(ce.Path, rpath)
			if _, err := createFolder(p); err != nil {
        // We should panic because this should not happen
        // and we can not recover from this
				log.Panic(err)
			}

		case EnvimLuaFile:
      rpath := path.Join(".nvim", "envim.lua")
      p := path.Join(ce.Path, rpath)
      if err := createConfigFile(p); err != nil {
        log.Panic(err)
      }
		case PluginsFolder:
			rpath := path.Join(".nvim", "plugins")
			p := path.Join(ce.Path, rpath)
			if _, err := createFolder(p); err != nil {
				log.Panic(err)
			}

		default:
			log.Panicf("%s: There is no implementation for %s EnvRepairFlag", re.GetType(), rflag)

		}
	}

	re.state = HandlerSuccessState
}
