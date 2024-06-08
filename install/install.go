package install

import "fmt"
import "envim/config"

func Install(file string, force bool) (bool, error) {
  configMap, err := config.ReadConfig(file)
  if err != nil {
    return false, err
  }
  fmt.Println(configMap)
  return true, nil
}
