package config

import (
	"errors"
	"github.com/yuin/gopher-lua"
)

var keys = [...]string{"dependencies", "nvim_version", "plugin_manager"}

func ReadConfig(file string) (map[string]lua.LValue, error) {

  // Make a map to hold the key values
  keyValues := make(map[string]lua.LValue)

  // Open a new lua state
	L := lua.NewState()
	defer L.Close()

  // Run dofile to push the file return value to the stack
	err := L.DoFile(file)
	if err != nil {
		return keyValues, err
	}

  // Get the return value from the stack
	lv := L.Get(-1)

  // Check if the return value is a table
	if val, ok := lv.(*lua.LTable); ok {

    // Check if the table has the required keys
		for _, key := range keys {
      if _val := val.RawGetString(key); _val == lua.LNil {
				return keyValues, errors.New("missing key: " + key)
			} else {
          keyValues[key] = _val
      }

		}
		return keyValues, nil
	} else {
		return keyValues, errors.New("envim file must return a table")
	}
}
