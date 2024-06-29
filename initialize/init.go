package initialize

import (
	"errors"
	"fmt"
	"os"
	"path"
)

var EnvimDir string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	EnvimDir = path.Join(homeDir, ".envim")

	if m, err := os.Stat(EnvimDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(EnvimDir, 0755); err != nil {
				panic(err)
			}
		} else {
			if !m.IsDir() {
				panic(errors.New(fmt.Sprintf("%s %s", EnvimDir, "exists but is not a directory")))
			}
		}
	}

	if m, err := os.Stat(path.Join(EnvimDir, "versions")); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(EnvimDir, "versions"), 0777); err != nil {
				panic(err)
			}
		} else {
			if !m.IsDir() {
				panic(errors.New(fmt.Sprintf("%s %s", path.Join(EnvimDir, "versions"), "exists but is not a directory")))
			}
		}
	}

}
