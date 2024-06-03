package main

import (
	"envim/validate"
	"fmt"
)

func main() {

  v, err := validate.Validate("envim.lua")
  if err != nil {
    panic(err)
  }
  fmt.Println(v)
}
