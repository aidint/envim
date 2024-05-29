#!/bin/bash

nvim_version=${1}
rm -rf ~/.local/share/envim/v${nvim_version}
echo "downloading neovim v${nvim_version}"
git clone --depth 1 --branch v${nvim_version} https://github.com/neovim/neovim.git ~/.local/share/envim/clone/v${nvim_version} > /dev/null 2>&1
echo "building neovim v${nvim_version}"
cd ~/.local/share/envim/clone/v${nvim_version}
make CMAKE_BUILD_TYPE=RelWithDebInfo CMAKE_INSTALL_PREFIX=~/.local/share/envim/v${nvim_version} > /dev/null 2>&1
make install > /dev/null 2>&1
echo "neovim v${nvim_version} installed"

