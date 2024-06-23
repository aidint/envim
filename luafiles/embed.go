package luafiles
import (
  _ "embed"
)

//go:embed files/sample_config.lua
var SampleConfig string
