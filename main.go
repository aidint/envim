package main

import (
	"envim/cmd"
)

import (
	"errors"
	"fmt"
	"os"
	"path"
)

func main() {
	cmd.Execute()
}


func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
  envimDir := path.Join(homeDir, ".envim")

	if m, err := os.Stat(envimDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(envimDir, 0777); err != nil {
				panic(err)
			}
		} else {
			if !m.IsDir() {
				panic(errors.New(fmt.Sprintf("%s %s", envimDir, "exists but is not a directory")))
			}
		}
	}

	if m, err := os.Stat(path.Join(envimDir, "versions")); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(envimDir, "versions"), 0777); err != nil {
				panic(err)
			}
		} else {
			if !m.IsDir() {
				panic(errors.New(fmt.Sprintf("%s %s", path.Join(envimDir, "versions"), "exists but is not a directory")))
			}
		}
	}

}
