local M = {}

M.nvim_version = os.getenv('NVIM_VERSION') or "v0.10.0"
M.plugin_manager = {name = 'folke/lazynvim', tag = 'v10.22.1'}

M.dependencies = {
  ['nvim-lua/plenary.nvim'] = os.getenv('NVIM_VERSION') or "v1.22";
  ['nvim-lua/popup.nvim'] = "v10.21.2";
  ['nvim-telescope/telescope.nvim'] = function (environ) return environ["somevalue"] end;
}

return M
