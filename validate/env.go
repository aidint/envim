package validate

import (
	"errors"
	"os"
)

// This function should later validate the whole environment including
// .envim folder structure
// .nvim folder structure
// validate the nvim version health
func ValidateEnvironment() error {
  if _, err := os.Stat(".envim"); err != nil {
    if os.IsNotExist(err) {
      return errors.New("Envim environment has not been initialized yet")
    }
    return err
  }
  return nil
}

