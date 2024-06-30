package run

import (
	"envim/initialize"
	"errors"
	"os"
	"os/exec"
	"path"
)

func Run(m map[string]interface{}, environment map[string]string, args []string) error {

  nvim_version, err := extractNvimVersion(m)
  if err != nil {
    return err
  }

  nvim := path.Join(initialize.EnvimDir, "versions", nvim_version, "envim", "bin/nvim")
  cmd := exec.Command(nvim, args...)
  cmd.Stdout = os.Stdout
  cmd.Stdin = os.Stdin
  cmd.Stderr = os.Stderr
  if err := cmd.Run(); err != nil {
    return err
  }
  return nil
}

func RunDefault(args []string) error {
  cmd := exec.Command("nvim", args...)
  cmd.Stdout = os.Stdout
  cmd.Stdin = os.Stdin
  cmd.Stderr = os.Stderr
  if err := cmd.Run(); err != nil {
    return err
  }
  return nil
}

func extractNvimVersion(m map[string]interface{}) (string, error) {
  if val, ok := m["nvim"]; ok {
    if val, ok := val.(map[string]interface{}); ok {
      for key, value := range val {
        if key == "version" {
          if nvimVersion, ok := value.(string); ok {
            return nvimVersion, nil
          } else {
            return "", errors.New("Error reading envim.json: nvim version is not a string")
          }
        }
      }
    } else {
      return "", errors.New("Error reading envim.json: nvim key is not a map")
    }
  } else {
    return "", errors.New("Error reading envim.json: nvim key not found")
  }
  return "", errors.New("nvim version not found in envim.json")
}

