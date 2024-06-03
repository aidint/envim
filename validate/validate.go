package validate

import (
	"errors"
	"fmt"

	"github.com/yuin/gopher-lua"
)

func Validate(filepath string) (bool, error) {

	L := lua.NewState()
	defer L.Close()

	err := L.DoFile(filepath)
	if err != nil {
		return false, err
	}

	lv := L.Get(-1)
	if val, ok := lv.(*lua.LTable); ok {
		fmt.Println(*val)
	} else {
		return false, errors.New("envim file must return a table")
	}

	return true, nil
}
