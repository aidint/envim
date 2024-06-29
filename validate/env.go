package validate

import (
	"errors"
	"os"
)

var EnvimExists bool
// This function should later validate the whole environment including
// .envim folder structure
// .nvim folder structure
// validate the nvim version health
func ValidateEnvironment() error {
  if !EnvimExists {
    return errors.New("Envim environment has not been initialized yet")
  }
  return nil
}

func init() {

  if res, err := os.Stat(".envim"); err != nil {
    EnvimExists = false
    return
  } else if !res.IsDir() {
    EnvimExists = false
    return
  }

  EnvimExists = true
}

