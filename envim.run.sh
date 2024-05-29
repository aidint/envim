#!/bin/bash

nvim_version=${1}
shift
exec ~/.local/share/envim/v${nvim_version}/bin/nvim $@
