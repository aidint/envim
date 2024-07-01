package config

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/yuin/gopher-lua"
)

var keys = [...]string{"dependencies", "nvim_version", "plugin_manager"}
type TableConfig map[string]interface{}

func (t TableConfig) GetValue(keys ...string) (interface{}, error) {
  var value interface{}
  value = t

  for _, key := range keys{
    if v, ok := value.(map[string]interface{})[key]; ok {
      value = v
    } else {
      return nil, errors.New(fmt.Sprintf("Key {%s} not found", key))
    }
  }
  return value, nil
}

func GetStringValue(L *lua.LState, val lua.LValue, args ...lua.LValue) (string, error) {
	switch v := val.(type) {
	case lua.LString:
		return v.String(), nil
	case lua.LNumber:
		return strconv.FormatFloat(float64(v), 'f', -1, 64), nil
	case *lua.LFunction:
		if err := L.CallByParam(lua.P{
			Fn:      v,
			NRet:    1,
			Protect: true,
		}, args...); err != nil {
			return "", errors.New(fmt.Sprintf("Lua function error -> %s", err))
		}
		lv := L.Get(-1)
		if s, ok := lv.(lua.LString); ok {
			return s.String(), nil
		} else {
			return "", errors.New("Lua function must return a string")
		}
	default:
		return "", errors.New("Value is neither a string nor a function")
	}
}

// A function that reads a table that stores only string-compatible values
// or tables with the same property. String-compatible values are: strings, numbers and functions that return strings
func ReadTable(L *lua.LState, table *lua.LTable, args ...lua.LValue) (map[string]interface{}, error) {

	keyValues := make(map[lua.LValue]lua.LValue)
	table.ForEach(func(key, val lua.LValue) {
		keyValues[key] = val
	})

	m := make(map[string]interface{})
	for k, v := range keyValues {
		var key string
    var value interface{}
		key, err := GetStringValue(L, k, args...)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error reading key -> %s", err))
		}
		switch val := v.Type(); val {
		case lua.LTTable:
      value, err = ReadTable(L, v.(*lua.LTable), args...)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("%s -> %s", key, err))
			}
    
		case lua.LTNumber, lua.LTString, lua.LTFunction:
      value, err = GetStringValue(L, v, args...)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Error reading value for key {%s} -> %s", key, err))
			}
		default:
			return nil, errors.New(fmt.Sprintf("{%s} -> Unsupported value type", key))
		}
    m[key] = value
	}

	return m, nil
}

func ReadConfig(L *lua.LState, file string, args ...lua.LValue) (map[string]lua.LValue, error) {

	// Make a map to hold the key values
	keyValues := make(map[string]lua.LValue)

	// Run dofile to push the file return value to the stack
	err := L.DoFile(file)
	if err != nil {
		return nil, err
	}

	// Get the return value from the stack
	lv := L.Get(-1)

	// Check if the return value is a table
	if val, ok := lv.(*lua.LTable); ok {

    table, err := ReadTable(L, val, args...)
    if err != nil {
      fmt.Println(err)
    }
    fmt.Printf("%+v\n", table)
		// Check if the table has the required keys
		// This should be later replaced by a generic validation function
		for _, key := range keys {
			if _val := val.RawGetString(key); _val == lua.LNil {
				return nil, errors.New("missing key {" + key + "}")
			} else {
				keyValues[key] = _val
				_val.Type()
			}
		}

		return keyValues, nil
	} else {
		return nil, errors.New("envim file must return a table")
	}
}
