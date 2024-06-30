package config

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/yuin/gopher-lua"
)

var keys = [...]string{"dependencies", "nvim_version", "plugin_manager"}

func ReadConfig(L *lua.LState, file string) (map[string]lua.LValue, error) {

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

		// Check if the table has the required keys
		// This should be later replaced by a generic validation function
		for _, key := range keys {
			if _val := val.RawGetString(key); _val == lua.LNil {
				return nil, errors.New("missing key: " + key)
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

func GetStringValue(L *lua.LState, val lua.LValue, args ...lua.LValue) (string, error) {
	switch v := val.(type) {
	case lua.LString:
		return v.String(), nil
  case lua.LNumber:
    return strconv.Itoa(int(v)), nil
	case *lua.LFunction:
		if err := L.CallByParam(lua.P{
			Fn:      v,
			NRet:    1,
			Protect: true,
		}, args...); err != nil {
			return "", errors.New(fmt.Sprintf("Lua function error: %s", err))
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
