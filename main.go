package main

import (
	"envim/validate"
)

func main() {

  _, err := validate.Validate("envim.lua")
  if err != nil {
    panic(err)
  }
}
