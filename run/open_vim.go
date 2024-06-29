package run

import (
	"envim/initialize"
	"os"
	"os/exec"
	"path"
)

func Run(nvim_version string, environment map[string]string) error {
  nvim := path.Join(initialize.EnvimDir, "versions", nvim_version, "envim", "bin/nvim")
  cmd := exec.Command(nvim)
  cmd.Stdout = os.Stdout
  cmd.Stdin = os.Stdin
  cmd.Stderr = os.Stderr
  if err := cmd.Run(); err != nil {
    return err
  }
  return nil
}
