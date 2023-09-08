function abbr(a, w)
	vim.cmd("iabbrev " .. a .. " " .. w)
end

vim.opt.expandtab = false
vim.opt.shiftwidth = 4
vim.opt.tabstop = 4
abbr("ee", ":=")
