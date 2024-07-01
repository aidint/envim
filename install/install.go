package install

import (
	"envim/initialize"
	"errors"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/yuin/gopher-lua"
)

func InstallNvim(nvimVersion string, force bool) error {

	dir := path.Join(initialize.EnvimDir, "versions", nvimVersion)
	if res, err := os.Stat(dir); err == nil {
		if !res.IsDir() {
			os.Remove(dir)
		} else if !force {
			log.Printf("Neovim version %s already installed. Skipping...\n", nvimVersion)
			return nil
		} else if force {
			log.Printf("Neovim version %s already installed. Forcing reinstallation...\n", nvimVersion)
			os.RemoveAll(dir)
		}
	}

	log.Printf("Downloading neovim version %s\n", nvimVersion)
	_, err := git.PlainClone(path.Join(initialize.EnvimDir, "versions", nvimVersion), false, &git.CloneOptions{
		URL:           "https://github.com/neovim/neovim",
		SingleBranch:  true,
		Depth:         1,
		ReferenceName: plumbing.ReferenceName("refs/tags/" + nvimVersion),
	})

	if err != nil {
		return err
	}

	log.Printf("Building neovim version %s\n", nvimVersion)
	e := exec.Command("make", "CMAKE_BUILD_TYPE=RelWithDebInfo", "CMAKE_INSTALL_PREFIX="+path.Join(initialize.EnvimDir, "versions", nvimVersion, "envim"))
	e.Dir = dir
	err = e.Run()
	if err != nil {
		return err
	}

	log.Printf("Installing neovim version %s\n", nvimVersion)
	e = exec.Command("make", "install")
	e.Dir = dir
	err = e.Run()
	if err != nil {
		return err
	}
	return nil
}

// Returns map of of installed dependencies and their versions
func Install(L *lua.LState, configMap map[string]interface{}, force bool) (map[string]interface{}, error) {

	m := make(map[string]interface{})

  var nvimVersion string
  nvimVersion, ok := configMap["nvim_version"].(string)


	if !ok {
    return nil, errors.New("nvim_version is not a string-compatible value") 
  }

  if err := InstallNvim(nvimVersion, force); err != nil {
    return nil, err
  }

  m["nvim"] = map[string]string{
      "version": nvimVersion,
  }

	return m, nil
}
