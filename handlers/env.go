package handlers

// import "os"

type EnvRepairFlag int
const (
  PluginsFolder = iota 
)

type CheckEnvironment struct {
  state HandlerState
  errors []error
  rflags []EnvRepairFlag
  Path string
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
  // if info, err := os.Stat(ce.Path); 
}
