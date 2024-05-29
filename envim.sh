#!/bin/bash

# thia script takes two arguments: the command and the version. based on command, it will do the following:
# - install: install the specified version of neovim
# - run   : run the specified version of neovim

command=${1}
if [ "${command}" == "init" ]; then
  if [ -d .envim ]; then
    echo "envim is already initialized"
  else
    echo "initializing envim"
    mkdir .envim
    touch .envim/init.lua
    echo "envim initialized"
  fi
fi

if [ -d .envim ]; then
  config_dir=$(pwd)
  app_name=".envim"
else
  echo "envim is not initialized"
  exit 1
fi

if [ "${command}" == "install" ]; then
  nvim_version=${2}
  bash ./envim.add.sh ${nvim_version}
elif [ "${command}" == "run" ]; then
  nvim_version=${2}
  XDG_CONFIG_HOME=${config_dir} NVIM_APPNAME=${app_name} bash ./envim.run.sh ${nvim_version}
fi
