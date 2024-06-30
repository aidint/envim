package run

import (
	"log"
	"os"
	"path"
)

var xdg map[string]string

func init() {

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	xdg = map[string]string{
		"XDG_CACHE_HOME":  path.Join(wd, ".envim", "cache"),
		"XDG_CONFIG_HOME": path.Join(wd, ".envim", "config"),
		"XDG_DATA_HOME":   path.Join(wd, ".envim", "data"),
		"XDG_RUNTIME_DIR": path.Join(wd, ".envim", "runtime_dir"),
		"XDG_CONFIG_DIRS": path.Join(wd, ".envim", "config_dirs"),
		"XDG_DATA_DIRS":   path.Join(wd, ".envim", "data_dirs"),
		"XDG_STATE_HOME":  path.Join(wd, ".envim", "state"),
	}
}

func exportEnv() []string {
	var env []string
	for key, value := range xdg {
		env = append(env, key+"="+value)
	}
	return env
}
