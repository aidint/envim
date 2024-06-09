package install

import (
	"envim/config"
	"errors"
	"os"
	"path"

	"github.com/yuin/gopher-lua"
)

var envimDir string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	envimDir = path.Join(homeDir, ".envim")
}

func installNvim(nvim_version string) error {
  
	return nil
}

func Install(file string, force bool) (bool, error) {
	configMap, err := config.ReadConfig(file)
	if err != nil {
		return false, err
	}

	if nvim_version, ok := configMap["nvim_version"].(lua.LString); ok {
		installNvim(nvim_version.String())
	} else {
		return false, errors.New("nvim_version must be a string")
	}

	return true, nil
}
