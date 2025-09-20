local markdown = require("markdown")

-- md renders a markdown string to an html string. This is mainly for convience in dealing with long form written posts.
-- They don't need to be part of the component rendering system
function md(s)
	return {__safe=true, markdown(s)}
end

-- safe performs some basic escaping of html values so that text can be put in an html page without needing to worry about parsing issues.
-- I'm sure I missed some values in here that should be escaped I'll add them at a later point
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

local selfClosing = {
  area = true,
  base = true,
  br = true,
  col = true,
  embed = true,
  hr = true,
  img = true,
  input = true,
  link = true,
  meta = true,
  param = true,
  source = true,
  track = true,
  wbr = true
}
-- render is the function that does all the work. It takes tables, and recursively uses the write function to render them. 
-- the basic structure of the 'node' table is `{__name="div", child, child, arg1="value" arg2="value2"}`
function render(node, write)
  -- if a name is set, then this is an html node, so lets render its opening tag with its parameters
  -- This also allow us to nest tables and then unest them: {{{"some"}, "text"}, "content"}
	if node.__name ~= nil then
    -- write the start of the node "<div"
		write("<",node.__name)
    -- NOTE: we use pairs to grab all the key=value members
    -- this lets us deal with only the attributes: {id="me"}
		for idx, value in pairs(node) do
			if type(idx) == "string" and idx ~= "__name" then
        -- we can accept a table of strings and write all of them out.
        -- this is to support setting multiple classes on a node:
        -- class={"first", "second"}
        -- This way we can dynamically add to the classes without needing to combine them into a single string
				if type(value) == "table" then
          -- " id="
					write(' ',idx,'="')
					local first = ""
					for _,v in ipairs(value) do
						write(first, safe(v))
						first = " "
					end
					write('"')
				else
          -- just a simple value, lets write it out
					write(' ',idx,'="',value,'"')
				end
			end
		end
    if selfClosing[node.__name] then
      write("/")
    end
		write(">")
	end
  -- self closing elements, or void elements, can not have a body
  if selfClosing[node.__name] then
    return
  end
  -- Now we have an opening tag, we can continue and write out the children on this node.
  -- NOTE: we are using ipairs to specifically loop over the table as if it was an array
  -- this allows us to only deal with the children: {"asdf","foo"}
	local first = ""
	for idx, child in ipairs(node) do
		if type(child) == "table" then
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

local oldRequire = require
require = function(path)
  return oldRequire("lua/endpoints/" .. path)
end

-- This sets up the global envrionment
local G = _G
local newGT = {}
setmetatable(newGT,{
	__index= function(_, n)
		return G[n] or function(params)
			params.__name=n
			return params
		end
	end,
})
setfenv(0, newGT)

