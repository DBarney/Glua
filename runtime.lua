local markdown = require("lua/markdown")

function md(s)
	return {__safe=true, markdown(s)}
end

local function safe(s)
	s = string.gsub(s, "[&<>'\"]", {
		["&"]="&amp;",
		["<"]="&lt;",
		[">"]="&gt;",
		["'"]="&#39;",
		["\""]="&quot;",
	})
	return s
end

function render(node, write)
	if node.__name ~= nil then
		write("<",node.__name)
		for idx, value in pairs(node) do
			if type(idx) == "string" and idx ~= "__name" then
				if type(value) == "table" then
						write(' ',idx,'="')
						local first = ""
					for _,v in ipairs(value) do
						write(first, safe(v))
						first = " "
					end
						write('"')
					else
					write(' ',idx,'="',value,'"')
				end
			end
		end
		write(">")
	end
	local first = ""
	for idx, child in ipairs(node) do
		if type(child) == "table" then
			write(first)
			render(child, write)
		elseif type(child) == "string" then
			if node.__safe then
				write(first, child)
			else
				write(first, safe(child))
			end
		end
		first = " "
	end
	if node.__name ~= nil then
		write("</",node.__name, ">")
	end
end

local G = _G
local newGT = {}
setmetatable(newGT,{
	__index= function(_, n)
		if not G[n] then
			return function(params)
				params.__name=n
				return params
			end
		else
			return G[n]
		end
	end,
})
setfenv(0, newGT)
