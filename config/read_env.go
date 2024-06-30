package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

func ReadEnv() (map[string]interface{}, error) {
  file, err := os.ReadFile(path.Join(".envim", "envim.json"))
  if err != nil {
    if os.IsNotExist(err) {
      return nil, errors.New("envim.json not found, try installing dependencies first by running `envim install`")
    }
    return nil, err
  }

  var env map[string]interface{}
  if err := json.Unmarshal(file, &env); err != nil {
    return nil, err
  }

  return env, nil
}
