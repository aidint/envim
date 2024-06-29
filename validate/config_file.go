package validate

import (
	"errors"
	"github.com/yuin/gopher-lua"
)

var keys = [...]string{"dependencies", "nvim_version", "plugin_manager"}

func Validate(filepath string) (bool, error) {

	L := lua.NewState()
	defer L.Close()

	err := L.DoFile(filepath)
	if err != nil {
		return false, err
	}

	lv := L.Get(-1)
	if val, ok := lv.(*lua.LTable); ok {
		for _, key := range keys {
			if val.RawGetString(key) == lua.LNil {
				return false, errors.New("missing key: " + key)
			}
		}
		return true, nil
	} else {
		return false, errors.New("envim file must return a table")
	}

}
