package initialize

import (
	"envim/luafiles"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
)

var workingDir string

func init() {
	workingDir, _ = os.Getwd()
}

type FlagData struct {
	Active bool
	Value  string
}

func CreateEnvironment() (string, error) {
	environmentPath := path.Join(workingDir, ".envim")

	if res, err := os.Stat(".envim"); err == nil {
		if res.IsDir() {
			return environmentPath, errors.New("Environment folder already exists in " + workingDir)
		} else {
			return environmentPath, errors.New("A file by the name '.envim' already exists in the current directory")
		}
	}

	if err := os.Mkdir(".envim", 0755); err != nil {
		return environmentPath, err
	}

  if err := createGitIgnore(); err != nil {
    return environmentPath, err
  }

	return environmentPath, nil
}

func CreateDotNvim() (string, error) {
	dotnvimPath := path.Join(workingDir, ".nvim")

	if res, err := os.Stat(".nvim"); err == nil {
		if res.IsDir() {
			return dotnvimPath, errors.New("'.nvim' folder already exists in " + workingDir)
		} else {
			return dotnvimPath, errors.New("A file by the name '.nvim' already exists in the current directory")
		}
	}

	if err := os.Mkdir(".nvim", 0755); err != nil {
		return dotnvimPath, err
	}
	return dotnvimPath, nil
}

func createGitIgnore() error {

	file, err := os.OpenFile(path.Join(".envim", ".gitignore"), os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString("*"); err != nil {
		return err
	}
	return nil
}

func CreateConfigFile(configFile string) (string, error) {

  var configFilePath string
  // Check if the path is absolute
  if configFilePath = path.Join(workingDir, configFile); path.IsAbs(configFile) {
    configFilePath = configFile
  }

	if _, err := os.Stat(configFile); err == nil {
		return configFilePath, errors.New(fmt.Sprintf("%s already exists in the current directory", configFile))
	}

	if err := os.WriteFile(configFile, []byte(luafiles.SampleConfig), 0644); err != nil {
		return configFilePath, err
	}
	log.Printf("%s file created in %s\n", configFile, workingDir)
	return configFilePath, nil
}
