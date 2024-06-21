package install

import (
	"envim/config"
	"errors"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/yuin/gopher-lua"
)

var envimDir string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not get user home directory")
	}
	envimDir = path.Join(homeDir, ".envim")
}

func installNvim(nvim_version string, force bool) error {

	dir := path.Join(envimDir, "versions", nvim_version)
	if res, err := os.Stat(dir); err == nil {
		if !res.IsDir() {
			os.Remove(dir)
		} else if !force {
			log.Printf("Neovim version %s already installed. Skipping...\n", nvim_version)
			return nil
		} else if force {
			log.Printf("Neovim version %s already installed. Forcing reinstallation...\n", nvim_version)
			os.RemoveAll(dir)
		}
	}

	log.Printf("Downloading neovim version %s\n", nvim_version)
	_, err := git.PlainClone(path.Join(envimDir, "versions", nvim_version), false, &git.CloneOptions{
		URL:           "https://github.com/neovim/neovim",
		SingleBranch:  true,
		Depth:         1,
		ReferenceName: plumbing.ReferenceName("refs/tags/" + nvim_version),
	})

	if err != nil {
		return err
	}

	log.Printf("Building neovim version %s\n", nvim_version)
	e := exec.Command("make", "CMAKE_BUILD_TYPE=RelWithDebInfo", "CMAKE_INSTALL_PREFIX="+path.Join(envimDir, "versions", nvim_version, "envim"))
	e.Dir = dir
	err = e.Run()
	if err != nil {
		return err
	}

	log.Printf("Installing neovim version %s\n", nvim_version)
	e = exec.Command("make", "install")
	e.Dir = dir
	err = e.Run()
	if err != nil {
		return err
	}
	return nil
}

func Install(file string, force bool) (bool, error) {
	configMap, err := config.ReadConfig(file)
	if err != nil {
		return false, err
	}

	if nvim_version, ok := configMap["nvim_version"].(lua.LString); ok {
		err = installNvim(nvim_version.String(), force)
	} else {
		return false, errors.New("nvim_version must be a string")
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
